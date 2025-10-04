package store

type StoreValueType interface {
	GetType() string
}

type RawBytes struct {
	B []byte
}

func (t RawBytes) GetType() string {
	return "string"
}

type List struct {
	L    []string
	Null bool
}

func (l List) IsEmpty() bool {
	return l.Null || len(l.L) == 0
}

func (t List) GetType() string {
	return "list"
}
