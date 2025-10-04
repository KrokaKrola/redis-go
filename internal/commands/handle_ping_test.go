package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// TestParse_PING_BulkString verifies that Parse accepts an array
// with a bulk-string command name "PING" and no arguments.
func TestParse_PING_BulkString(t *testing.T) {
	input := &resp.Array{Elems: []resp.Value{&resp.BulkString{B: []byte("PING")}}}

	cmd, perr := Parse(input)
	if perr != nil {
		t.Fatalf("Parse returned protocol error: %#v", perr)
	}
	if cmd == nil {
		t.Fatalf("Parse returned nil command")
	}
	if cmd.Name != PING_COMMAND {
		t.Fatalf("unexpected command name: got %q, want %q", cmd.Name, PING_COMMAND)
	}
	if len(cmd.Args) != 0 {
		t.Fatalf("unexpected args length: got %d, want 0", len(cmd.Args))
	}
}

// TestParse_PING_SimpleString verifies that a simple-string command name also works.
func TestParse_PING_SimpleString(t *testing.T) {
	input := &resp.Array{Elems: []resp.Value{&resp.SimpleString{S: []byte("PING")}}}

	cmd, perr := Parse(input)
	if perr != nil {
		t.Fatalf("Parse returned protocol error: %#v", perr)
	}
	if cmd == nil {
		t.Fatalf("Parse returned nil command")
	}
	if cmd.Name != PING_COMMAND {
		t.Fatalf("unexpected command name: got %q, want %q", cmd.Name, PING_COMMAND)
	}
	if len(cmd.Args) != 0 {
		t.Fatalf("unexpected args length: got %d, want 0", len(cmd.Args))
	}
}

func TestDispatch_PING_WithArg_ReturnsBulkString(t *testing.T) {
	cmd := &Command{
		Name: PING_COMMAND,
		Args: []resp.Value{&resp.BulkString{B: []byte("hello")}},
	}

	out := Dispatch(cmd, store.NewStore())
	bs, ok := out.(*resp.BulkString)
	if !ok {
		t.Fatalf("expected BulkString, got %T", out)
	}
	if bs.Null {
		t.Fatalf("unexpected null BulkString response")
	}
	if string(bs.B) != "hello" {
		t.Fatalf("unexpected PING response: got %q, want %q", string(bs.B), "hello")
	}
}

// TestDispatch_PING_TooManyArgs ensures that PING with more than one argument
// fails with an error.
func TestDispatch_PING_TooManyArgs(t *testing.T) {
	cmd := &Command{
		Name: PING_COMMAND,
		Args: []resp.Value{
			&resp.SimpleString{S: []byte("one")},
			&resp.SimpleString{S: []byte("two")},
		},
	}

	out := Dispatch(cmd, store.NewStore())
	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for too many args, got %T", out)
	}
}
