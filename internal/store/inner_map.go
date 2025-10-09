package store

import (
	"fmt"
	"time"
)

type innerMap map[string]storeValue

type storeValue struct {
	value      StoreValueType
	expiryTime time.Time
}

func (v storeValue) isExpired() bool {
	timeDiff := v.expiryTime.Compare(time.Now())

	return timeDiff <= 0
}

func newStoreValue(value StoreValueType, expiryTime time.Time) storeValue {
	return storeValue{
		value,
		expiryTime,
	}
}

func (m innerMap) xrange(key, start, end string) (Stream, error) {
	v, ok := m[key]

	if !ok {
		return Stream{}, nil
	}

	if v.isExpired() {
		return Stream{}, nil
	}

	stream, ok := v.value.(Stream)
	if !ok {
		return Stream{}, fmt.Errorf("MISSTYPE of the element in the underlying array")
	}

	// todo: add filtering by start and end

	return stream, nil
}
