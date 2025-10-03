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

func (l List) IsEmpty() bool {
	return l.Null || len(l.L) == 0
}

func (t List) isValue() {}
