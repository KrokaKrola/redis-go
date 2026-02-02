package resp

import (
	"bufio"
)

type CountingReader struct {
	R     *bufio.Reader
	Count int
}

func (c *CountingReader) Read(p []byte) (n int, err error) {
	n, err = c.R.Read(p)
	c.Count += n
	return n, err
}

func (c *CountingReader) ReadByte() (byte, error) {
	b, err := c.R.ReadByte()
	if err == nil {
		c.Count += 1
	}
	return b, err
}

func (c *CountingReader) Peek(n int) ([]byte, error) {
	return c.R.Peek(n)
}

func (c *CountingReader) ReadString(delim byte) (string, error) {
	s, err := c.R.ReadString(delim)
	c.Count += len(s)
	return s, err
}
