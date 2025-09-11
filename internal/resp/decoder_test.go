package resp

import (
	"bufio"
	"bytes"
	"testing"
)

// TestDecoder_Read_PINGArray verifies that the decoder parses a single-element
// array containing a bulk string "PING" from the RESP input.
func TestDecoder_Read_PINGArray(t *testing.T) {
	input := "*1\r\n$4\r\nPING\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)
	v, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() returned error: %v", err)
	}

	var arr *Array
	switch a := v.(type) {
	case *Array:
		arr = a
	default:
		t.Fatalf("expected Array, got %T", v)
	}

	if arr == nil {
		t.Fatalf("expected not-nil array")
	}

	if arr.Null {
		t.Fatalf("expected non-null array")
	}
	if len(arr.Elems) != 1 {
		t.Fatalf("expected array length 1, got %d", len(arr.Elems))
	}

	// Validate the single element is BulkString("PING")
	elem := arr.Elems[0]
	switch bs := elem.(type) {
	case *BulkString:
		if bs.Null {
			t.Fatalf("expected non-null bulk string")
		}
		if string(bs.B) != "PING" {
			t.Fatalf("expected bulk string 'PING', got %q", string(bs.B))
		}
	default:
		t.Fatalf("expected BulkString element, got %T", elem)
	}
}

// TestDecoder_Read_ECHOArray verifies that the decoder parses a two-element
// array ["ECHO", "hey"] from the RESP input.
func TestDecoder_Read_ECHOArray(t *testing.T) {
	input := "*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)
	v, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() returned error: %v", err)
	}

	var arr *Array
	switch a := v.(type) {
	case *Array:
		arr = a
	default:
		t.Fatalf("expected Array, got %T", v)
	}

	if arr == nil {
		t.Fatalf("expected not-nil array")
	}
	if arr.Null {
		t.Fatalf("expected non-null array")
	}
	if got := len(arr.Elems); got != 2 {
		t.Fatalf("expected array length 2, got %d", got)
	}

	// First element: command name "ECHO"
	name := arr.Elems[0]
	switch n := name.(type) {
	case *BulkString:
		if n.Null {
			t.Fatalf("expected non-null bulk string for name")
		}
		if string(n.B) != "ECHO" {
			t.Fatalf("expected command name 'ECHO', got %q", string(n.B))
		}
	case *SimpleString:
		if string(n.S) != "ECHO" {
			t.Fatalf("expected command name 'ECHO', got %q", string(n.S))
		}
	default:
		t.Fatalf("expected string-like name, got %T", name)
	}

	// Second element: argument "hey"
	arg := arr.Elems[1]
	switch bs := arg.(type) {
	case *BulkString:
		if bs.Null {
			t.Fatalf("expected non-null bulk string for arg")
		}
		if string(bs.B) != "hey" {
			t.Fatalf("expected arg 'hey', got %q", string(bs.B))
		}
	case *SimpleString:
		if string(bs.S) != "hey" {
			t.Fatalf("expected arg 'hey', got %q", string(bs.S))
		}
	default:
		t.Fatalf("expected BulkString/SimpleString arg, got %T", arg)
	}
}
