package store

import (
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
