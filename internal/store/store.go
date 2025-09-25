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

func (m innerMap) get(key string) (value []byte, ok bool, expired bool) {
	v, ok := m[key]

	if ok {
		timeDiff := v.expiryTime.Compare(time.Now())

		if timeDiff <= 0 {
			return nil, true, true
		}
	}

	return v.value, ok, false
}

func (m innerMap) set(key string, value []byte, expType ExpiryType, expiryTime int) bool {
	t := time.Now()

	switch expType {
	case EXPIRY_EX:
		t = t.Add(time.Duration(expiryTime) * time.Second)
	case EXPIRY_PX:
		t = t.Add(time.Duration(expiryTime) * time.Millisecond)
	case "":
		t = time.Date(9999, 12, 31, 23, 59, 59, 999, time.UTC)
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

	// allow empty string for indefinite key storage
	if vLower == "" {
		return "", true
	}

	if vLower != string(EXPIRY_EX) && vLower != string(EXPIRY_PX) {
		return "", false
	}

	return ExpiryType(vLower), true
}
