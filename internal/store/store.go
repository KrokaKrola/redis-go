package store

import (
	"fmt"
	"sync"
	"time"
)

type Store struct {
	sync.RWMutex
	innerMap
	blpopQueue []string
}

func NewStore() *Store {
	return &Store{innerMap: make(innerMap)}
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
	return s.append(key, cpy)
}

func (s *Store) Lpush(key string, value []string) (int64, bool) {
	s.Lock()
	defer s.Unlock()
	cpy := append([]string{}, value...)
	return s.prepend(key, cpy)
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
	blpopId := newID()

	s.Lock()
	s.blpopQueue = append(s.blpopQueue, blpopId)
	s.Unlock()

	for {
		s.RLock()
		l, ok, expired, wrongType := s.getList(key)

		if wrongType {
			s.RUnlock()
			return "", false
		}

		if !ok || expired {
			// wait until list will be created
			s.RUnlock()
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if l.Null || len(l.L) == 0 {
			s.RUnlock()
			time.Sleep(10 * time.Millisecond)
			continue
		}

		if len(s.blpopQueue) > 0 && s.blpopQueue[0] != blpopId {
			// item is waited by another client
			s.RUnlock()
			time.Sleep(10 * time.Millisecond)
			continue
		}

		s.RUnlock()
		s.Lock()
		defer s.Unlock()

		list, ok := s.lpop(key, 1)
		fmt.Printf("for me %#v, ok: %t...\n", list, ok)
		if !ok {
			return "", false
		}

		s.blpopQueue = s.blpopQueue[1:]

		// not sure if it's needed since it will be checked before
		// if list.Null || len(list.L) == 0 {
		// 	time.Sleep(10 * time.Millisecond)
		// 	continue
		// }

		return list.L[0], true
	}
}
