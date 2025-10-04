package commands

import (
	"bufio"
	"bytes"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// TestParse_SET_BulkStrings verifies that Parse accepts an array
// with bulk-string command name "SET" and two bulk-string arguments.
func TestParse_SET_BulkStrings(t *testing.T) {
	input := &resp.Array{Elems: []resp.Value{
		&resp.BulkString{B: []byte("SET")},
		&resp.BulkString{B: []byte("mykey")},
		&resp.BulkString{B: []byte("myval")},
	}}

	cmd, perr := Parse(input)
	if perr != nil {
		t.Fatalf("Parse returned protocol error: %#v", perr)
	}
	if cmd == nil {
		t.Fatalf("Parse returned nil command")
	}
	if cmd.Name != SET_COMMAND {
		t.Fatalf("unexpected command name: got %q, want %q", cmd.Name, SET_COMMAND)
	}
	if len(cmd.Args) != 2 {
		t.Fatalf("unexpected args length: got %d, want 2", len(cmd.Args))
	}
	// Validate args are preserved in order with expected values.
	if bs, ok := cmd.Args[0].(*resp.BulkString); !ok || string(bs.B) != "mykey" {
		t.Fatalf("unexpected first arg: %#v", cmd.Args[0])
	}
	if bs, ok := cmd.Args[1].(*resp.BulkString); !ok || string(bs.B) != "myval" {
		t.Fatalf("unexpected second arg: %#v", cmd.Args[1])
	}
}

// TestParse_GET_BulkString verifies that Parse accepts an array
// with bulk-string command name "GET" and one bulk-string argument.
func TestParse_GET_BulkString(t *testing.T) {
	input := &resp.Array{Elems: []resp.Value{
		&resp.BulkString{B: []byte("GET")},
		&resp.BulkString{B: []byte("mykey")},
	}}

	cmd, perr := Parse(input)
	if perr != nil {
		t.Fatalf("Parse returned protocol error: %#v", perr)
	}
	if cmd == nil {
		t.Fatalf("Parse returned nil command")
	}
	if cmd.Name != GET_COMMAND {
		t.Fatalf("unexpected command name: got %q, want %q", cmd.Name, GET_COMMAND)
	}
	if len(cmd.Args) != 1 {
		t.Fatalf("unexpected args length: got %d, want 1", len(cmd.Args))
	}
	if bs, ok := cmd.Args[0].(*resp.BulkString); !ok || string(bs.B) != "mykey" {
		t.Fatalf("unexpected GET key arg: %#v", cmd.Args[0])
	}
}

// TestDispatch_SET_Then_GET_ReturnsValue verifies that using SET to store a value
// and then GET on the same key returns that value as a bulk string.
func TestDispatch_SET_Then_GET_ReturnsValue(t *testing.T) {
	s := store.NewStore()

	// SET mykey myval
	setCmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
			&resp.BulkString{B: []byte("myval")},
			&resp.BulkString{B: []byte("px")},
			&resp.BulkString{B: []byte("100")},
		},
	}
	out1 := Dispatch(setCmd, s)
	if ss, ok := out1.(*resp.SimpleString); !ok || string(ss.S) != "OK" {
		t.Fatalf("expected +OK for SET, got %#v", out1)
	}

	// GET mykey
	getCmd := &Command{
		Name: GET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
		},
	}
	out2 := Dispatch(getCmd, s)
	bs, ok := out2.(*resp.BulkString)
	if !ok {
		t.Fatalf("expected BulkString for GET, got %T", out2)
	}
	if bs.Null {
		t.Fatalf("expected non-null BulkString for existing key")
	}
	if string(bs.B) != "myval" {
		t.Fatalf("unexpected value: got %q, want %q", string(bs.B), "myval")
	}
}

// TestDispatch_SET_WithPXExpiry ensures that SET with PX expiry stores the value
// temporarily and GET returns null after the expiry window.
func TestDispatch_SET_WithPXExpiry(t *testing.T) {
	s := store.NewStore()

	raw := "*5\r\n$3\r\nSET\r\n$5\r\napple\r\n$5\r\nmango\r\n$2\r\npx\r\n$3\r\n100\r\n"
	dec := resp.NewDecoder(bufio.NewReader(bytes.NewReader([]byte(raw))))

	value, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read returned error: %v", err)
	}

	cmd, perr := Parse(value)
	if perr != nil {
		t.Fatalf("Parse returned protocol error: %#v", perr)
	}

	out := Dispatch(cmd, s)
	if ss, ok := out.(*resp.SimpleString); !ok || string(ss.S) != "OK" {
		t.Fatalf("expected +OK for SET with PX, got %#v", out)
	}

	getCmd := &Command{
		Name: GET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("apple")},
		},
	}

	immediate := Dispatch(getCmd, s)
	bs, ok := immediate.(*resp.BulkString)
	if !ok || bs.Null {
		t.Fatalf("expected BulkString for GET before expiry, got %#v", immediate)
	}
	if string(bs.B) != "mango" {
		t.Fatalf("unexpected value before expiry: got %q, want %q", string(bs.B), "mango")
	}

	time.Sleep(150 * time.Millisecond)

	expired := Dispatch(getCmd, s)
	if bs, ok := expired.(*resp.BulkString); !ok || !bs.Null {
		t.Fatalf("expected null BulkString after expiry, got %#v", expired)
	}
}

// TestDispatch_GET_Nonexistent_ReturnsNull verifies that GET on a missing key
// returns a null bulk string.
func TestDispatch_GET_Nonexistent_ReturnsNull(t *testing.T) {
	s := store.NewStore()
	getCmd := &Command{
		Name: GET_COMMAND,
		Args: []resp.Value{&resp.BulkString{B: []byte("missing")}},
	}
	out := Dispatch(getCmd, s)
	bs, ok := out.(*resp.BulkString)
	if !ok {
		t.Fatalf("expected BulkString, got %T", out)
	}
	if !bs.Null {
		t.Fatalf("expected null BulkString for missing key")
	}
}
