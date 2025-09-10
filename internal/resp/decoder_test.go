package resp

import (
	"bufio"
	"bytes"
	"fmt"
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
		fmt.Println("Array???")
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
