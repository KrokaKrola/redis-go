package resp

import (
	"bufio"
	"fmt"
)

type Decoder struct {
	r *bufio.Reader
}

func NewDecoder(reader *bufio.Reader) *Decoder {
	return &Decoder{
		r: reader,
	}
}

func (d *Decoder) Read() (Value, error) {
	b, err := d.r.ReadByte()

	if err != nil {
		return nil, err
	}

	switch b {
	case byte('+'):
		value, err := d.processSimpleString()

		if err != nil {
			return nil, err
		}

		return value, nil
	case byte('*'):
		value, err := d.processArray()

		if err != nil {
			return nil, err
		}

		return value, nil
	case byte('$'):
		value, err := d.processBulkString()

		if err != nil {
			return nil, err
		}

		return value, nil
	default:
		return nil, fmt.Errorf("unknown data type: %s", string(b))
	}
}

func (d *Decoder) processSimpleString() (*SimpleString, error) {
	return nil, nil
}

func (d *Decoder) processArray() (*Array, error) {
	return nil, nil
}

func (d *Decoder) processBulkString() (*BulkString, error) {
	return nil, nil
}
