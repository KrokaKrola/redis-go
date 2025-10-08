package store

import "time"

func (m innerMap) set(key string, value []byte, expType ExpiryType, expiryTime int) bool {
	t, ok := getExpiryTime(expType, expiryTime)

	if !ok {
		return false
	}

	m[key] = newStoreValue(RawBytes{Bytes: value}, t)
	return true
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
