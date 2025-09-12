package store

import "sync"

type innerMap map[string][]byte

func (m innerMap) get(key string) ([]byte, bool) {
	v, ok := m[key]
	return v, ok
}

func (m innerMap) set(key string, value []byte) {
	m[key] = value
}

func (m innerMap) delete(key string) {
	delete(m, key)
}

type Store struct {
	sync.RWMutex
	innerMap
}

func NewStore() *Store {
	return &Store{innerMap: make(innerMap)}
}

func (s *Store) Get(key string) ([]byte, bool) {
	s.RLock()
	defer s.RUnlock()
	value, ok := s.get(key)

	if !ok {
		return nil, ok
	}

	cValue := make([]byte, len(value))
	copy(cValue, value)
	return cValue, ok
}

func (s *Store) Set(key string, value []byte) {
	s.Lock()
	defer s.Unlock()
	cValue := make([]byte, len(value))
	copy(cValue, value)
	s.set(key, cValue)
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	s.delete(key)
}
