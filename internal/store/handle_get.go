package store

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
