package store

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

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

type innerMap map[string]storeValue

func (m innerMap) get(key string) (value []byte, ok bool, expired bool) {
	v, ok := m[key]

	if ok {
		if v.isExpired() {
			return nil, true, true
		}
	}

	rb, isRawBytes := v.value.(RawBytes)
	if !isRawBytes {
		return nil, false, false
	}

	return rb.Bytes, ok, false
}

func getPossibleEndTime() time.Time {
	return time.Date(9999, 12, 31, 23, 59, 59, 999, time.UTC)
}

func getExpiryTime(expType ExpiryType, expTime int) (keyExpTime time.Time, ok bool) {
	t := time.Now()

	switch expType {
	case EXPIRY_EX:
		t = t.Add(time.Duration(expTime) * time.Second)
	case EXPIRY_PX:
		t = t.Add(time.Duration(expTime) * time.Millisecond)
	case "":
		t = getPossibleEndTime()
	default:
		return t, false
	}

	return t, true
}

func (m innerMap) set(key string, value []byte, expType ExpiryType, expiryTime int) bool {
	t, ok := getExpiryTime(expType, expiryTime)

	if !ok {
		return false
	}

	m[key] = newStoreValue(RawBytes{Bytes: value}, t)
	return true
}

func (m innerMap) append(key string, arr []string) (int64, bool) {
	v, ok := m[key]

	if !ok || v.isExpired() {
		m[key] = newStoreValue(List{Elements: arr}, getPossibleEndTime())
		return int64(len(arr)), true
	}

	list, isList := v.value.(List)
	if !isList {
		return 0, false
	}

	newArr := append(list.Elements, arr...)
	m[key] = newStoreValue(List{Elements: newArr}, v.expiryTime)
	return int64(len(newArr)), true
}

func (m innerMap) prepend(key string, arr []string) (int64, bool) {
	v, ok := m[key]

	if !ok || v.isExpired() {
		slices.Reverse(arr)
		m[key] = newStoreValue(List{Elements: arr}, getPossibleEndTime())
		return int64(len(arr)), true
	}

	list, isList := v.value.(List)
	if !isList {
		return 0, false
	}

	slices.Reverse(arr)
	arr = append(arr, list.Elements...)
	m[key] = newStoreValue(List{Elements: arr}, v.expiryTime)
	return int64(len(arr)), true
}

func (m innerMap) delete(key string) {
	delete(m, key)
}

func (m innerMap) getList(key string) (list List, ok bool, expired bool, wrongType bool) {
	v, ok := m[key]

	if !ok {
		return List{Null: true}, false, false, false
	}

	if v.isExpired() {
		return List{Null: true}, true, true, false
	}

	l, isList := v.value.(List)
	if !isList {
		return List{Null: true}, false, false, true
	}

	return l, true, false, false
}

func (m innerMap) lpop(key string, count int) (list List, ok bool) {
	v, ok := m[key]

	if !ok {
		return List{Null: true}, true
	}

	if v.isExpired() {
		return List{Null: true}, true
	}

	l, isList := v.value.(List)

	if !isList {
		return List{Null: true}, false
	}

	if l.Null || len(l.Elements) == 0 {
		return List{Null: true}, true
	}

	listLen := len(l.Elements)

	if count > listLen {
		count = listLen
	}

	elems := l.Elements[:count]
	newList := l.Elements[count:]

	m[key] = newStoreValue(List{Elements: newList}, v.expiryTime)

	return List{Elements: elems}, true
}

func (m innerMap) getRawValue(key string) (value StoreValueType, ok bool) {
	sv, ok := m[key]

	if !ok {
		return nil, false
	}

	if sv.isExpired() {
		delete(m, key)
		return nil, false
	}

	return sv.value, true
}

func (m innerMap) xadd(key string, msTime uint64, seqNumber uint64, fields [][]string) (newEntryId string, err error) {
	sv, ok := m[key]

	fields = cloneStreamFields(fields)

	id := fmt.Sprintf("%d-%d", msTime, seqNumber)

	if !ok || sv.isExpired() {
		m[key] = newStoreValue(Stream{
			Elements:           []streamElement{{id, fields}},
			LtsInsertedIdParts: streamIdParts{msTime, seqNumber},
		}, getPossibleEndTime())

		return id, nil
	}

	stream, okStream := sv.value.(Stream)
	if !okStream {
		return "", errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
	}

	if isValidStreamId := validateStreamIdParts(msTime, seqNumber, stream.LtsInsertedIdParts); !isValidStreamId {
		return "", fmt.Errorf("ERR The ID specified in XADD is equal or smaller than the target stream top item")
	}

	newElements := append(stream.Elements, streamElement{id, fields})

	m[key] = newStoreValue(Stream{
		Elements:           newElements,
		LtsInsertedIdParts: streamIdParts{msTime, seqNumber},
	}, sv.expiryTime)

	return id, nil
}

func cloneStreamFields(fields [][]string) [][]string {
	if len(fields) == 0 {
		return nil
	}

	out := make([][]string, len(fields))
	for i, pair := range fields {
		copied := make([]string, len(pair))
		copy(copied, pair)
		out[i] = copied
	}

	return out
}

func validateStreamIdParts(msTime uint64, seqNumber uint64, ltsInsertedIdParts streamIdParts) bool {
	if msTime < ltsInsertedIdParts.msTime {
		return false
	}

	if msTime == ltsInsertedIdParts.msTime {
		if seqNumber <= ltsInsertedIdParts.seqNumber {
			return false
		}
	}

	return true
}
