package store

import "time"

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

	return rb.B, ok, false
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

	m[key] = newStoreValue(RawBytes{B: value}, t)
	return true
}

func (m innerMap) append(key string, arr []string) (int64, bool) {
	v, ok := m[key]

	if !ok || v.isExpired() {
		m[key] = newStoreValue(List{L: arr}, getPossibleEndTime())
		return int64(len(arr)), true
	}

	switch it := v.value.(type) {
	case List:
		newArr := append(it.L, arr...)
		m[key] = newStoreValue(List{L: newArr}, v.expiryTime)
		return int64(len(newArr)), true
	default:
		return 0, false
	}
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
