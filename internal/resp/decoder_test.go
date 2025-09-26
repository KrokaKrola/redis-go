package resp

import (
	"bufio"
	"bytes"
	"fmt"
	"math"
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

// TestDecoder_Read_NullBulkStringThenSimpleString ensures that a null bulk string
// ("$-1\r\n") is parsed and the following value is read correctly (consumes CRLF).
func TestDecoder_Read_NullBulkStringThenSimpleString(t *testing.T) {
	input := "$-1\r\n+OK\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)

	v1, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() v1 returned error: %v", err)
	}
	bs, ok := v1.(*BulkString)
	if !ok {
		t.Fatalf("expected first value BulkString, got %T", v1)
	}
	if !bs.Null {
		t.Fatalf("expected null bulk string for first value")
	}

	v2, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() v2 returned error: %v", err)
	}
	ss, ok := v2.(*SimpleString)
	if !ok {
		t.Fatalf("expected second value SimpleString, got %T", v2)
	}
	if string(ss.S) != "OK" {
		t.Fatalf("expected simple string 'OK', got %q", string(ss.S))
	}
}

// TestDecoder_Read_BulkStringEmpty validates that an empty bulk string ($0) is
// decoded as a non-null bulk string with zero-length payload.
func TestDecoder_Read_BulkStringEmpty(t *testing.T) {
	input := "$0\r\n\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)
	v, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() returned error: %v", err)
	}

	bs, ok := v.(*BulkString)
	if !ok {
		t.Fatalf("expected BulkString, got %T", v)
	}
	if bs.Null {
		t.Fatalf("expected non-null bulk string for $0")
	}
	if len(bs.B) != 0 {
		t.Fatalf("expected empty payload for $0, got %d bytes", len(bs.B))
	}
}

// TestDecoder_Read_Integer parses positive integer value as an
// Integer type
func TestDecoder_Read_Integer(t *testing.T) {
	input := ":5\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)
	v, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() returned error: %v", err)
	}

	i, ok := v.(*Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", v)
	}

	if i.N != 5 {
		t.Fatalf("expected integer with value 5, got value %d", i.N)
	}
}

// TestDecoder_Read_Integer_WithOptionalPlusSign parses positive integer value
// with optional plus sign as an Integer value
func TestDecoder_Read_Integer_WithOptionalPlusSign(t *testing.T) {
	input := ":+5\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)
	v, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() returned error: %v", err)
	}

	i, ok := v.(*Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", v)
	}

	if i.N != 5 {
		t.Fatalf("expected integer with value 5, got value %d", i.N)
	}
}

// TestDecoder_Read_Integer_WithNegativeSign parses negative integer value
// with negative sign as an Integer value
func TestDecoder_Read_Integer_WithNegativeSign(t *testing.T) {
	input := ":-5\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)
	v, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() returned error: %v", err)
	}

	i, ok := v.(*Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", v)
	}

	if i.N != -5 {
		t.Fatalf("expected integer with value 5 and IsNegative = true, got value %d", i.N)
	}
}

// TestDecoder_Read_IntegerFollowedBySimpleString verifies that decoding an
// integer consumes its CRLF and the subsequent frame can be read normally.
func TestDecoder_Read_IntegerFollowedBySimpleString(t *testing.T) {
	input := ":1\r\n+OK\r\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)

	v1, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() v1 returned error: %v", err)
	}
	i, ok := v1.(*Integer)
	if !ok {
		t.Fatalf("expected first value Integer, got %T", v1)
	}
	if i.N != 1 {
		t.Fatalf("expected integer 1, got %d", i.N)
	}

	v2, err := dec.Read()
	if err != nil {
		t.Fatalf("decoder.Read() v2 returned error: %v", err)
	}
	ss, ok := v2.(*SimpleString)
	if !ok {
		t.Fatalf("expected second value SimpleString, got %T", v2)
	}
	if string(ss.S) != "OK" {
		t.Fatalf("expected simple string 'OK', got %q", string(ss.S))
	}
}

// TestDecoder_Read_IntegerRejectsLFOnly ensures RESP integers terminated with
// LF-only are rejected since RESP requires CRLF.
func TestDecoder_Read_IntegerRejectsLFOnly(t *testing.T) {
	input := ":1\n"
	br := bufio.NewReader(bytes.NewBufferString(input))

	dec := NewDecoder(br)

	if _, err := dec.Read(); err == nil {
		t.Fatalf("expected decoder.Read() to fail for LF-only terminator")
	}
}

// TestRESP_IntegerExtremesRoundTrip validates decoding and re-encoding of
// RESP integer boundary values.
func TestRESP_IntegerExtremesRoundTrip(t *testing.T) {
	cases := []int64{math.MaxInt64, math.MinInt64}

	for _, tc := range cases {
		input := fmt.Sprintf(":%d\r\n", tc)
		br := bufio.NewReader(bytes.NewBufferString(input))

		dec := NewDecoder(br)
		v, err := dec.Read()
		if err != nil {
			t.Fatalf("decoder.Read() returned error for %d: %v", tc, err)
		}

		i, ok := v.(*Integer)
		if !ok {
			t.Fatalf("expected Integer for %d, got %T", tc, v)
		}
		if i.N != tc {
			t.Fatalf("expected integer %d, got %d", tc, i.N)
		}

		var buf bytes.Buffer
		enc := NewEncoder(&buf)
		if err := enc.Write(i); err != nil {
			t.Fatalf("encoder.Write() returned error for %d: %v", tc, err)
		}

		if got := buf.String(); got != input {
			t.Fatalf("unexpected round-trip output for %d: got %q, want %q", tc, got, input)
		}
	}
}
