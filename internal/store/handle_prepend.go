package store

import "slices"

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
