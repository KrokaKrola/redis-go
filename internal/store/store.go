package store

import (
	"sync"
)

type Store struct {
	sync.RWMutex
	innerMap
	blpopQueue map[string][]chan string
}

func NewStore() *Store {
	return &Store{
		innerMap:   make(innerMap),
		blpopQueue: make(map[string][]chan string),
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
		return List{Null: true}, false
	}

	if !ok {
		s.RUnlock()
		return List{Null: true}, true
	}

	if expired {
		s.RUnlock()

		s.Lock()
		defer s.Unlock()

		s.delete(key)
		return List{Null: true}, true
	}

	if storeList.Null {
		s.RUnlock()
		return List{Null: true}, true
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
		return List{Null: true}, true
	}

	if start >= storeListLen {
		s.RUnlock()
		return List{Null: true}, true
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

func (s *Store) Blpop(key string, timeoutInSeconds int) (el string, ok bool) {
	res, ok := s.lpop(key, 1)

	if ok && !res.Null && len(res.L) > 0 {
		return res.L[0], true
	}

	listener := make(chan string, 1)
	s.blpopQueue[key] = append(s.blpopQueue[key], listener)
	el = <-listener

	return el, true
}

func (s *Store) produceElementToListeners(key string) {
	queue, ok := s.blpopQueue[key]

	if !ok || len(queue) == 0 {
		return
	}

	for len(queue) > 0 {
		res, ok := s.lpop(key, 1)

		if !ok || res.Null || len(res.L) == 0 {
			break
		}

		queue[0] <- res.L[0]
		queue = queue[1:]
	}

	if len(queue) == 0 {
		delete(s.blpopQueue, key)
		return
	}

	s.blpopQueue[key] = queue
}
