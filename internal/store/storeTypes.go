package store

type StoreValueType interface {
	isValue()
}

type RawBytes struct {
	B []byte
}

func (t RawBytes) isValue() {}

type List struct {
	L    []string
	Null bool
}

func (t List) isValue() {}
