package resp

type Value interface {
	isValue()
}

type SimpleString struct {
	S []byte
}

func (s *SimpleString) isValue() {}

type BulkString struct {
	B    []byte
	Null bool
}

func (s *BulkString) isValue() {}

type Integer struct {
	N          int64
	IsNegative bool
}

func (s *Integer) isValue() {}

type Error struct {
	Msg string
}

func (s *Error) isValue() {}

type Array struct {
	Elems []Value
	Null  bool
}

func (s *Array) isValue() {}
