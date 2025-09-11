package resp

import (
	"bufio"
	"fmt"
	"strconv"
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

func (d *Decoder) getSizeOfTheData() (int, error) {
	var sizeStr []byte

	for {
		pBytes, err := d.r.Peek(1)

		if err != nil {
			return 0, err
		}

		if pBytes[0] == '\r' {
			break
		}

		rByte, err := d.r.ReadByte()

		if err != nil {
			return 0, err
		}

		sizeStr = append(sizeStr, rByte)
	}

	num, err := strconv.Atoi(string(sizeStr))

	if err != nil {
		return 0, err
	}

	return num, nil
}

func (d *Decoder) readClrf() error {
	_, err := d.r.Discard(2)

	if err != nil {
		return err
	}

	return nil
}

func (d *Decoder) processSimpleString() (*SimpleString, error) {
	str := &SimpleString{}

	for {
		pBytes, err := d.r.Peek(1)

		if err != nil {
			return nil, err
		}

		if pBytes[0] == '\r' {
			break
		}

		rByte, err := d.r.ReadByte()

		if err != nil {
			return nil, err
		}

		str.S = append(str.S, rByte)
	}

	if err := d.readClrf(); err != nil {
		return nil, err
	}

	return str, nil
}

func (d *Decoder) processArray() (*Array, error) {
	arrSize, err := d.getSizeOfTheData()

	if err != nil {
		return nil, err
	}

	if err := d.readClrf(); err != nil {
		return nil, err
	}

	arr := &Array{}

	if arrSize == 0 {
		arr.Null = true
		return arr, nil
	}

	for arrSize > len(arr.Elems) {
		el, err := d.Read()
		if err != nil {
			return nil, err
		}

		arr.Elems = append(arr.Elems, el)
	}

	return arr, nil
}

func (d *Decoder) processBulkString() (*BulkString, error) {
	strSize, err := d.getSizeOfTheData()

	if err != nil {
		return nil, err
	}

	if err := d.readClrf(); err != nil {
		return nil, err
	}

	str := &BulkString{}

	for len(str.B) < strSize {
		b, err := d.r.ReadByte()
		if err != nil {
			return nil, err
		}

		str.B = append(str.B, b)
	}

	if err := d.readClrf(); err != nil {
		return nil, err
	}

	return str, nil
}
