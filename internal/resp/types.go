package resp

type Value interface {
	isValue()
}

type SimpleString struct {
	Bytes []byte
}

func (s *SimpleString) isValue() {}

type BulkString struct {
	Bytes []byte
	Null  bool
}

func (s *BulkString) isValue() {}

type Integer struct {
	Number int64
}

func (s *Integer) isValue() {}

type Error struct {
	Msg string
}

func (s *Error) isValue() {}

type Array struct {
	Elements []Value
	Null     bool
}

func (s *Array) isValue() {}
