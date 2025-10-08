package store

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
