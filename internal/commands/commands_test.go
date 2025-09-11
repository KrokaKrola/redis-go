package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

// TestParse_PING_BulkString verifies that Parse accepts an array
// with a bulk-string command name "PING" and no arguments.
func TestParse_PING_BulkString(t *testing.T) {
	// *1 \r\n $4 \r\n PING \r\n in structured form
	input := resp.Array{Elems: []resp.Value{resp.BulkString{B: []byte("PING")}}}

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
	input := resp.Array{Elems: []resp.Value{resp.BulkString{B: []byte("PING")}}}

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
