package store

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	sync.RWMutex
	innerMap
	blpopQueue map[string][]blpopListener
	xreadQueue map[string][]xreadListener
}

type blpopListener struct {
	id      string
	valueCh chan string
}

type xreadEvent struct {
	element StreamElement
	isMax   bool
	key     string
}

type xreadListener struct {
	notify chan<- xreadEvent
	id     StreamIdSpec
}

func NewStore() *Store {
	return &Store{
		innerMap:   make(innerMap),
		blpopQueue: make(map[string][]blpopListener),
		xreadQueue: make(map[string][]xreadListener),
	}
}

func (s *Store) Get(key string) (value []byte, ok bool) {
	// lock read
	s.RLock()
	value, ok, expired := s.get(key)

	if !ok {
		// unlocking read when key is not found
		s.RUnlock()
		return nil, ok
	}

	if !expired {
		cpy := append([]byte{}, value...)
		s.RUnlock()
		return cpy, true
	}

	// unlock read if there is a key with expired time value
	s.RUnlock()
	// lock read & write since we will need to delete found value
	s.Lock()
	defer s.Unlock()

	value, ok, expired = s.get(key)

	if !ok || expired {
		if ok && expired {
			s.delete(key)
		}

		return nil, false
	}

	cpy := append([]byte{}, value...)
	return cpy, true
}

func (s *Store) Set(key string, value []byte, expType ExpiryType, expiryTime int) bool {
	s.Lock()
	defer s.Unlock()
	cpy := append([]byte{}, value...)
	return s.set(key, cpy, expType, expiryTime)
}

func (s *Store) Rpush(key string, value []string) (int64, bool) {
	s.Lock()
	defer s.Unlock()

	cpy := append([]string{}, value...)
	len, ok := s.append(key, cpy)

	if ok {
		s.produceElementToListeners(key)
	}

	return len, ok
}

func (s *Store) Lpush(key string, value []string) (int64, bool) {
	s.Lock()
	defer s.Unlock()

	cpy := append([]string{}, value...)
	len, ok := s.prepend(key, cpy)

	if ok {
		s.produceElementToListeners(key)
	}

	return len, ok
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	s.delete(key)
}

func (s *Store) Lrange(key string, start int, stop int) (list List, ok bool) {
	s.RLock()

	storeList, ok, expired, wrongType := s.getList(key)

	if wrongType {
		s.RUnlock()
		return List{}, false
	}

	if !ok {
		s.RUnlock()
		return List{}, true
	}

	if expired {
		s.RUnlock()

		s.Lock()
		defer s.Unlock()

		s.delete(key)
		return List{}, true
	}

	if storeList.Null {
		s.RUnlock()
		return List{}, true
	}

	storeListLen := len(storeList.Elements)

	if start < 0 {
		start = max(storeListLen+start, 0)
	}

	if stop < 0 {
		stop = max(storeListLen+stop, -1)
	}

	if start > stop {
		s.RUnlock()
		return List{}, true
	}

	if start >= storeListLen {
		s.RUnlock()
		return List{}, true
	}

	stop += 1

	if stop > storeListLen {
		stop = storeListLen
	}

	s.RUnlock()

	return List{Elements: storeList.Elements[start:stop]}, true
}

func (s *Store) Lpop(key string, count int) (list List, ok bool) {
	s.Lock()
	defer s.Unlock()
	return s.lpop(key, count)
}

func (s *Store) Blpop(key string, timeoutInSeconds float64) (el string, ok bool, timeout bool) {
	s.Lock()

	res, ok := s.lpop(key, 1)

	if !ok {
		s.Unlock()
		return "", false, false
	}

	if ok && !res.IsEmpty() {
		s.Unlock()
		return res.Elements[0], true, false
	}

	listener := blpopListener{
		id:      newId(),
		valueCh: make(chan string, 1),
	}
	s.blpopQueue[key] = append(s.blpopQueue[key], listener)

	s.Unlock()

	timeoutCh := (<-chan time.Time)(nil)
	if timeoutInSeconds > 0 {
		timeoutCh = time.After(time.Duration(timeoutInSeconds * float64(time.Second)))
	}

	select {
	case el := <-listener.valueCh:
		return el, true, false
	case <-timeoutCh:
		s.Lock()
		var newQueue []blpopListener

		for _, v := range s.blpopQueue[key] {
			if v.id != listener.id {
				newQueue = append(newQueue, v)
			}
		}

		s.blpopQueue[key] = newQueue

		s.Unlock()

		return "", false, true
	}

}

func (s *Store) produceElementToListeners(key string) {
	queue, ok := s.blpopQueue[key]

	if !ok || len(queue) == 0 {
		return
	}

	for len(queue) > 0 {
		res, ok := s.lpop(key, 1)

		if !ok || res.IsEmpty() {
			break
		}

		queue[0].valueCh <- res.Elements[0]
		queue = queue[1:]
	}

	if len(queue) == 0 {
		delete(s.blpopQueue, key)
		return
	}

	s.blpopQueue[key] = queue
}

func (s *Store) GetStoreRawValue(key string) (StoreValueType, bool) {
	s.Lock()
	defer s.Unlock()

	return s.getRawValue(key)
}

func (s *Store) Xadd(key string, streamId StreamIdSpec, fields [][]string) (newEntryId string, err error) {
	s.Lock()

	streamElement, err := s.xadd(key, streamId, fields)

	s.Unlock()

	if err != nil {
		return "", err
	}

	s.produceXaddEvents(key, streamElement)

	return fmt.Sprintf("%d-%d", streamElement.Id.MsTime, streamElement.Id.Seq), nil
}

func (s *Store) produceXaddEvents(key string, newElement StreamElement) {
	s.Lock()
	defer s.Unlock()

	listeners := s.xreadQueue[key]

	if len(listeners) > 0 {
		remaining := []xreadListener{}
		for _, listener := range listeners {
			event := xreadEvent{element: newElement, key: key}

			if listener.id.IsMax {
				event.isMax = true
				listener.notify <- event
				continue
			}

			if newElement.Id.MsTime < listener.id.MsTime {
				remaining = append(remaining, listener)
				continue
			}

			if newElement.Id.MsTime == listener.id.MsTime && newElement.Id.Seq <= listener.id.Seq {
				remaining = append(remaining, listener)
				continue
			}

			listener.notify <- event
		}

		if len(remaining) == 0 {
			delete(s.xreadQueue, key)
		} else {
			s.xreadQueue[key] = remaining
		}
	}
}

func (s *Store) Xrange(key string, start string, end string) (Stream, error) {
	s.Lock()
	defer s.Unlock()

	return s.xrange(key, start, end)
}

func (s *Store) Xread(keys [][]string, timeoutMs int, isBlocking bool) ([]Stream, error) {
	s.Lock()

	streams, err := s.xread(keys, false)

	if err != nil {
		s.Unlock()
		return nil, err
	}

	if !isBlocking {
		s.Unlock()
		return streams, nil
	}

	hasEntries := false

	for _, stream := range streams {
		if len(stream.Elements) > 0 {
			hasEntries = true
			break
		}
	}

	if hasEntries {
		s.Unlock()
		return streams, nil
	}

	notifyCh := make(chan xreadEvent, len(keys))

	cleanup := func() {
		for _, pair := range keys {
			queue := s.xreadQueue[pair[0]]
			if len(queue) == 0 {
				continue
			}

			filtered := queue[:0]
			for _, listener := range queue {
				if listener.notify != notifyCh {
					filtered = append(filtered, listener)
				}
			}

			if len(filtered) == 0 {
				delete(s.xreadQueue, pair[0])
			} else {
				s.xreadQueue[pair[0]] = filtered
			}
		}
	}

	for _, pair := range keys {
		id, ok := parseXreadStreamId(pair[1])
		if !ok {
			cleanup()
			s.Unlock()
			return []Stream{}, fmt.Errorf("ERR invalid stream id %s", pair[1])
		}

		s.xreadQueue[pair[0]] = append(s.xreadQueue[pair[0]], xreadListener{notify: notifyCh, id: id})
	}

	s.Unlock()

	timeoutCh := (<-chan time.Time)(nil)
	if timeoutMs > 0 {
		timeoutCh = time.After(time.Duration(float64(timeoutMs) * float64(time.Millisecond)))
	}

	select {
	case event := <-notifyCh:
		s.Lock()

		if event.isMax {
			cleanup()
			s.Unlock()
			streams := []Stream{}

			for _, pair := range keys {
				if pair[0] == event.key {
					streams = append(streams, Stream{Elements: []StreamElement{event.element}})
				} else {
					streams = append(streams, Stream{})
				}
			}

			return streams, nil
		}

		streams, err := s.xread(keys, true)

		cleanup()
		s.Unlock()

		return streams, err
	case <-timeoutCh:
		// didn't receive any elements during timeout
		s.Lock()
		cleanup()
		s.Unlock()

		return []Stream{}, nil
	}
}
