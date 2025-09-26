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
	case *BulkString:
		if v.Null {
			if _, err := fmt.Fprintf(e.w, "$-1\r\n"); err != nil {
				return err
			}

		} else {
			if _, err := fmt.Fprintf(e.w, "$%d\r\n%s\r\n", len(v.B), v.B); err != nil {
				return err
			}
		}
	case *Integer:
		if _, err := fmt.Fprintf(e.w, ":%d\r\n", v.N); err != nil {
			return err
		}
	case *Error:
		if _, err := fmt.Fprintf(e.w, "-%s\r\n", v.Msg); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown value type %T", v)
	}

	return nil
}
