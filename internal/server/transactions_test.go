package server

import (
	"bufio"
	"net"
	"sync"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func TestSingleConnectionTransactionFlow(t *testing.T) {
	srv := NewRedisServer(0, false)
	client := newInMemoryClient(t, srv)
	t.Cleanup(client.Close)

	requireSimpleString(t, client.do("MULTI"), "OK")
	requireSimpleString(t, client.do("SET", "foo", "bar"), "QUEUED")
	requireSimpleString(t, client.do("GET", "foo"), "QUEUED")

	execReply := requireArrayLen(t, client.do("EXEC"), 2)
	requireSimpleString(t, execReply.Elements[0], "OK")
	requireBulkString(t, execReply.Elements[1], "bar")

	requireBulkString(t, client.do("GET", "foo"), "bar")
}

func TestMultipleConnectionsTransactionFlow(t *testing.T) {
	srv := NewRedisServer(0, false)

	client1 := newInMemoryClient(t, srv)
	t.Cleanup(client1.Close)

	client2 := newInMemoryClient(t, srv)
	t.Cleanup(client2.Close)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		requireSimpleString(t, client1.do("MULTI"), "OK")
		requireSimpleString(t, client1.do("SET", "user:1", "alice"), "QUEUED")
		requireSimpleString(t, client1.do("GET", "user:1"), "QUEUED")

		execReply := requireArrayLen(t, client1.do("EXEC"), 2)
		requireSimpleString(t, execReply.Elements[0], "OK")
		requireBulkString(t, execReply.Elements[1], "alice")
	}()

	go func() {
		defer wg.Done()
		requireSimpleString(t, client2.do("MULTI"), "OK")
		requireSimpleString(t, client2.do("SET", "user:2", "bob"), "QUEUED")
		requireSimpleString(t, client2.do("GET", "user:2"), "QUEUED")

		execReply := requireArrayLen(t, client2.do("EXEC"), 2)
		requireSimpleString(t, execReply.Elements[0], "OK")
		requireBulkString(t, execReply.Elements[1], "bob")
	}()

	wg.Wait()

	verifier := newInMemoryClient(t, srv)
	t.Cleanup(verifier.Close)

	requireBulkString(t, verifier.do("GET", "user:1"), "alice")
	requireBulkString(t, verifier.do("GET", "user:2"), "bob")
}

type testClient struct {
	t      *testing.T
	conn   net.Conn
	writer *bufio.Writer
	enc    *resp.Encoder
	dec    *resp.Decoder
}

func newInMemoryClient(t *testing.T, srv *RedisServer) *testClient {
	t.Helper()

	clientConn, serverConn := net.Pipe()
	go srv.handleConnection(serverConn)

	writer := bufio.NewWriter(clientConn)

	return &testClient{
		t:      t,
		conn:   clientConn,
		writer: writer,
		enc:    resp.NewEncoder(writer),
		dec:    resp.NewDecoder(bufio.NewReader(clientConn)),
	}
}

func (c *testClient) Close() {
	c.conn.Close()
}

func (c *testClient) do(parts ...string) resp.Value {
	c.t.Helper()

	arr := &resp.Array{
		Elements: make([]resp.Value, len(parts)),
	}

	for i, part := range parts {
		arr.Elements[i] = &resp.BulkString{Bytes: []byte(part)}
	}

	if err := c.enc.Write(arr); err != nil {
		c.t.Fatalf("write %v: %v", parts, err)
	}

	if err := c.writer.Flush(); err != nil {
		c.t.Fatalf("flush %v: %v", parts, err)
	}

	reply, err := c.dec.Read()
	if err != nil {
		c.t.Fatalf("read %v: %v", parts, err)
	}

	return reply
}

func requireSimpleString(t *testing.T, v resp.Value, expected string) {
	t.Helper()

	str, ok := v.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected SimpleString %q, got %T", expected, v)
	}

	if string(str.Bytes) != expected {
		t.Fatalf("expected SimpleString %q, got %q", expected, string(str.Bytes))
	}
}

func requireBulkString(t *testing.T, v resp.Value, expected string) {
	t.Helper()

	bs, ok := v.(*resp.BulkString)
	if !ok || bs.Null {
		t.Fatalf("expected BulkString %q, got %#v", expected, v)
	}

	if string(bs.Bytes) != expected {
		t.Fatalf("expected BulkString %q, got %q", expected, string(bs.Bytes))
	}
}

func requireArrayLen(t *testing.T, v resp.Value, expected int) *resp.Array {
	t.Helper()

	arr, ok := v.(*resp.Array)
	if !ok {
		t.Fatalf("expected Array with len %d, got %T", expected, v)
	}

	if len(arr.Elements) != expected {
		t.Fatalf("expected Array len %d, got %d", expected, len(arr.Elements))
	}

	return arr
}
