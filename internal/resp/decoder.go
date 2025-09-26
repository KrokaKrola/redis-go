package resp

import (
	"bufio"
	"bytes"
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
	case byte(':'):
		value, err := d.processInteger()

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
		return 0, fmt.Errorf("ERR invalid size")
	}

	return num, nil
}

func (d *Decoder) readClrf() error {
	b := make([]byte, 2)
	n, err := d.r.Read(b)

	if err != nil {
		return err
	}

	if n != 2 {
		return fmt.Errorf("ERR invalid clrf delimiter")
	}

	if !bytes.Equal(b, []byte("\r\n")) {
		return fmt.Errorf("ERR invalid clrf delimiter")
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

	if arrSize < 0 {
		return &Array{
			Null: true,
		}, nil
	}

	arr := &Array{}

	if arrSize == 0 {
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

	if strSize == -1 {
		return &BulkString{
			Null: true,
		}, nil
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

func (d *Decoder) processInteger() (*Integer, error) {
	sign, err := d.r.Peek(1)

	if err != nil {
		return nil, err
	}

	var str []byte

	if sign[0] == '+' || sign[0] == '-' {
		d.r.ReadByte()

		if sign[0] == '-' {
			str = append(str, '-')
		}
	}

	integer := &Integer{}

	for {
		c, err := d.r.ReadByte()

		if err != nil {
			return nil, err
		}

		isDigit := c >= '0' && c <= '9'

		if !isDigit {
			return nil, fmt.Errorf("ERR invalid integer value")
		}

		str = append(str, c)

		nextChar, err := d.r.Peek(1)

		if err != nil {
			return nil, err
		}

		if nextChar[0] == '\r' {
			break
		}
	}

	if len(str) == 0 {
		return nil, fmt.Errorf("ERR invalid integer value")
	}

	res, err := strconv.ParseInt(string(str), 10, 64)
	if err != nil {
		return nil, err
	}

	integer.N = res

	if err := d.readClrf(); err != nil {
		return nil, err
	}

	return integer, nil
}
