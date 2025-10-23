package resp

import (
	"bytes"
	"testing"
)

// TestEncoder_Write_NullBulkString ensures that a null bulk string is encoded
// exactly as "$-1\r\n" with no extra data.
func TestEncoder_Write_NullBulkString(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&BulkString{Null: true}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := "$-1\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

// TestEncoder_Write_BulkStringEmpty ensures that an empty non-null bulk string
// is encoded as "$0\r\n\r\n".
func TestEncoder_Write_BulkStringEmpty(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&BulkString{Bytes: []byte("")}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := "$0\r\n\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

// TestEncoder_Write_PositiveInteger ensuges that a positive integer
// is encoded as ":5\r\n"
func TestEncoder_Write_PositiveInteger(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&Integer{Number: 5}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := ":5\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

// TestEncoder_Write_NegativeInteger ensuges that a negative integer
// is encoded as ":-5\r\n"
func TestEncoder_Write_NegativeInteger(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&Integer{Number: -5}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := ":-5\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

// TestEncoder_Write_Array ensuges that an array
// is encoded as a valid value
func TestEncoder_Write_Array(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&Array{
		Elements: []Value{
			&BulkString{Bytes: []byte("a")},
			&BulkString{Bytes: []byte("ab")},
			&BulkString{Bytes: []byte("abc")},
		},
	}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := "*3\r\n$1\r\na\r\n$2\r\nab\r\n$3\r\nabc\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

// TestEncoder_Write_Array_Empty ensuges that an array
// is encoded as a valid value if array is empty
func TestEncoder_Write_Array_Empty(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&Array{
		Elements: []Value{},
		Null:     false,
	}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := "*0\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

// TestEncoder_Write_Array_Empty ensuges that an array
// is encoded as a valid value if array is empty
func TestEncoder_Write_Array_Null(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&Array{
		Elements: []Value{},
		Null:     true,
	}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := "*-1\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}

func TestEncoder_Write_Inner_Arrays(t *testing.T) {
	var buf bytes.Buffer
	enc := NewEncoder(&buf)

	if err := enc.Write(&Array{
		Elements: []Value{
			&Array{
				Elements: []Value{
					&BulkString{Bytes: []byte("1-0")},
					&Array{
						Elements: []Value{
							&BulkString{Bytes: []byte("temperature")},
							&BulkString{Bytes: []byte("36")},
						},
					},
				},
			},
			&Array{
				Elements: []Value{
					&BulkString{Bytes: []byte("2-0")},
					&Array{
						Elements: []Value{
							&BulkString{Bytes: []byte("temperature")},
							&BulkString{Bytes: []byte("35")},
						},
					},
				},
			},
		},
	}); err != nil {
		t.Fatalf("encoder.Write() returned error: %v", err)
	}

	got := buf.String()
	want := "*2\r\n*2\r\n$3\r\n1-0\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n36\r\n*2\r\n$3\r\n2-0\r\n*2\r\n$11\r\ntemperature\r\n$2\r\n35\r\n"
	if got != want {
		t.Fatalf("unexpected output: got %q, want %q", got, want)
	}
}
