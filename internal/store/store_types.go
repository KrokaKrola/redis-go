package store

type StoreValueType interface {
	GetType() string
}

type RawBytes struct {
	Bytes []byte
}

func (t RawBytes) GetType() string {
	return "string"
}

type List struct {
	Elements []string
	Null     bool
}

func (l List) IsEmpty() bool {
	return l.Null || len(l.Elements) == 0
}

func (l List) GetType() string {
	return "list"
}

type Stream struct {
	Elements []streamElement
}

type streamElement struct {
	id     string
	fields [][]string
}

func (s Stream) GetType() string {
	return "stream"
}
