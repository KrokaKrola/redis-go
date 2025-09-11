package resp

import (
	"fmt"
	"io"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		w: writer,
	}
}

func (e *Encoder) Write(v Value) error {
	switch v := v.(type) {
	case SimpleString:
		fmt.Fprintf(e.w, "+%s\r\n", v.S)
		return nil
	default:
		return fmt.Errorf("unknown value type")
	}
}
