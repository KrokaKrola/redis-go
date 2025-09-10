package resp

import "io"

type Encoder struct {
	w io.Writer
}

func NewEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		w: writer,
	}
}

func (e *Encoder) Write(v Value) error {
	return nil
}
