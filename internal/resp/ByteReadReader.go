package resp

type ByteReadReader interface {
	Read(p []byte) (n int, err error)
	ReadByte() (byte, error)
	Peek(n int) ([]byte, error)
	ReadString(delim byte) (string, error)
}
