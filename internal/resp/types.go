package resp

import (
	"fmt"
	"strings"
)

type Value interface {
	isValue()
	String() string
}

type SimpleString struct {
	Bytes []byte
}

func (s *SimpleString) isValue() {}

func (s *SimpleString) String() string {
	return string(s.Bytes)
}

type BulkString struct {
	Bytes []byte
	Null  bool
}

func (s *BulkString) isValue() {}

func (s *BulkString) String() string {
	return string(s.Bytes)
}

type Integer struct {
	Number int64
}

func (s *Integer) isValue() {}

func (s *Integer) String() string {
	return fmt.Sprint(s.Number)
}

type Error struct {
	Msg string
}

func (s *Error) isValue() {}

func (s *Error) String() string {
	return s.Msg
}

type Array struct {
	Elements []Value
	Null     bool
}

func (s *Array) isValue() {}

func (s *Array) String() string {
	if s.Null {
		return "[]"
	}

	result := make([]string, 0, len(s.Elements))

	for _, v := range s.Elements {
		result = append(result, v.String())
	}

	return "[" + strings.Join(result, ",") + "]"
}
