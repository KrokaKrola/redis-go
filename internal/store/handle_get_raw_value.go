package store

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
