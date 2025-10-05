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
		if _, err := fmt.Fprintf(e.w, "+%s\r\n", v.Bytes); err != nil {
			return err
		}
	case *BulkString:
		if v.Null {
			if _, err := fmt.Fprintf(e.w, "$-1\r\n"); err != nil {
				return err
			}

		} else {
			if _, err := fmt.Fprintf(e.w, "$%d\r\n%s\r\n", len(v.Bytes), v.Bytes); err != nil {
				return err
			}
		}
	case *Integer:
		if _, err := fmt.Fprintf(e.w, ":%d\r\n", v.Number); err != nil {
			return err
		}
	case *Error:
		if _, err := fmt.Fprintf(e.w, "-%s\r\n", v.Msg); err != nil {
			return err
		}
	case *Array:
		if v.Null {
			if _, err := fmt.Fprintf(e.w, "*-1\r\n"); err != nil {
				return err
			}

			break
		}

		if len(v.Elements) == 0 {
			if _, err := fmt.Fprintf(e.w, "*0\r\n"); err != nil {
				return err
			}

			break
		}

		if _, err := fmt.Fprintf(e.w, "*%d\r\n", len(v.Elements)); err != nil {
			return err
		}

		for _, v := range v.Elements {
			bs, ok := v.(*BulkString)
			if !ok {
				return fmt.Errorf("MISSTYPE of the element in the underlying array")
			}

			if _, err := fmt.Fprintf(e.w, "$%d\r\n%s\r\n", len(bs.Bytes), bs.Bytes); err != nil {
				return err
			}
		}
	default:
		return fmt.Errorf("unknown value type %T", v)
	}

	return nil
}
