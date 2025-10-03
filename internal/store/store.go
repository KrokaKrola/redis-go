package store

import (
	"sync"
	"time"
)

type Store struct {
	sync.RWMutex
	innerMap
	blpopQueue map[string][]blpopListener
}

type blpopListener struct {
	id      string
	valueCh chan string
}

func NewStore() *Store {
	return &Store{
		innerMap:   make(innerMap),
		blpopQueue: make(map[string][]blpopListener),
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

	storeListLen := len(storeList.L)

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

	return List{L: storeList.L[start:stop]}, true
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
		return res.L[0], true, false
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

		queue[0].valueCh <- res.L[0]
		queue = queue[1:]
	}

	if len(queue) == 0 {
		delete(s.blpopQueue, key)
		return
	}

	s.blpopQueue[key] = queue
}
