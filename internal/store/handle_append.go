package store

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
