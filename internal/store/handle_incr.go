package store

import (
	"fmt"
	"strconv"
)

func (m innerMap) incr(key string) (list int64, err error) {
	valueRaw, ok := m[key]

	if !ok {
		m[key] = newStoreValue(RawBytes{Bytes: []byte("1")}, getPossibleEndTime())
		return 1, nil
	}

	value, ok := valueRaw.value.(RawBytes)
	if !ok {
		return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	valueInt, err := strconv.ParseInt(string(value.Bytes), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	valueInt += 1
	m[key] = newStoreValue(RawBytes{Bytes: fmt.Append(nil, valueInt)}, getPossibleEndTime())

	return valueInt, nil
}
