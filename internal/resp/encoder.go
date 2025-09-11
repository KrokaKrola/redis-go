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
	case *SimpleString:
		if _, err := fmt.Fprintf(e.w, "+%s\r\n", v.S); err != nil {
			return err
		}
		return nil
	case *BulkString:
		if _, err := fmt.Fprintf(e.w, "$%d\r\n%s\r\n", len(v.B), v.B); err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("unknown value type")
	}
}
