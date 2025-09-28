package store

import (
	"sync"
)

type Store struct {
	sync.RWMutex
	innerMap
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

func (s *Store) Rpush(key string, value List) (int64, bool) {
	s.Lock()
	defer s.Unlock()
	cpy := append([]string{}, value.L...)
	return s.append(key, cpy)
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	s.delete(key)
}

func (s *Store) Lrange(key string, start int, stop int) (list List, ok bool) {
	s.RLock()

	if start > stop {
		s.RUnlock()
		return List{Null: true}, true
	}

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

	if start >= storeListLen {
		s.RUnlock()
		return List{Null: true}, true
	}

	stop += 1

	if stop >= storeListLen {
		stop = storeListLen
	}

	s.RUnlock()

	return List{L: storeList.L[start:stop]}, true
}
