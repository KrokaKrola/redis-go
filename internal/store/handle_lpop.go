package store

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
