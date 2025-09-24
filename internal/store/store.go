package store

import (
	"strings"
	"sync"
	"time"
)

type ExpiryType string

const (
	EXPIRY_PX ExpiryType = "px"
	EXPIRY_EX ExpiryType = "ex"
)

type storeValue struct {
	value      []byte
	expiryTime time.Time
}

func newStoreValue(value []byte, expiryTime time.Time) storeValue {
	return storeValue{
		value,
		expiryTime,
	}
}

type innerMap map[string]storeValue

func (m innerMap) get(key string) ([]byte, bool) {
	v, ok := m[key]

	if ok {
		timeDiff := v.expiryTime.Compare(time.Now())

		if timeDiff <= 0 {
			m.delete(key)
			return nil, false
		}
	}

	return v.value, ok
}

func (m innerMap) set(key string, value []byte, expType string, expiryTime int) bool {
	expTypeValue, ok := ProcessExpType(expType)

	if !ok {
		return false
	}

	t := time.Now()

	switch expTypeValue {
	case EXPIRY_EX:
		t = t.Add(time.Duration(expiryTime) * time.Second)
	case EXPIRY_PX:
		t = t.Add(time.Duration(expiryTime) * time.Millisecond)
	default:
		return false
	}

	m[key] = newStoreValue(value, t)
	return true
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

func (s *Store) Set(key string, value []byte, expType string, expiryTime int) bool {
	s.Lock()
	defer s.Unlock()
	cValue := make([]byte, len(value))
	copy(cValue, value)
	return s.set(key, cValue, expType, expiryTime)
}

func (s *Store) Delete(key string) {
	s.Lock()
	defer s.Unlock()
	s.delete(key)
}

func ProcessExpType(v string) (ExpiryType, bool) {
	vLower := strings.ToLower(v)

	if vLower != string(EXPIRY_EX) && vLower != string(EXPIRY_PX) {
		return "", true
	}

	return ExpiryType(vLower), false
}
