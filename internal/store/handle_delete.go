package store

func (m innerMap) delete(key string) {
	delete(m, key)
}
