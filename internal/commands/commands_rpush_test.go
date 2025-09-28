package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// TestDispatch_RPUSH confirms RPUSH appends one element and reports the new length.
func TestDispatch_RPUSH(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("foo")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)
	in, ok := out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.N != 1 {
		t.Fatalf("unexpected RPUSH response, got %d, want %d", in.N, 1)
	}

	cmd = &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("bar")},
		},
	}

	out = Dispatch(cmd, store)
	in, ok = out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.N != 2 {
		t.Fatalf("unexpected RPUSH response, got %d, want %d", in.N, 2)
	}
}

// TestDispatch_RPUSH_MultipleElements covers appending multiple elements in one call and the reported length.
func TestDispatch_RPUSH_MultipleElements(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("foo")},
			&resp.BulkString{B: []byte("bar")},
			&resp.BulkString{B: []byte("baz")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)
	in, ok := out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.N != 3 {
		t.Fatalf("unexpected RPUSH response, got %d, want %d", in.N, 3)
	}

	cmd = &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("xyz")},
			&resp.BulkString{B: []byte("qwerty")},
		},
	}

	out = Dispatch(cmd, store)
	in, ok = out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.N != 5 {
		t.Fatalf("unexpected RPUSH response, got %d, want %d", in.N, 5)
	}
}

// TestDispatch_RPUSH_EmptyListError verifies RPUSH errors when no values are provided.
func TestDispatch_RPUSH_EmptyListError(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)
	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for invalid amount of arguments, got %T", out)
	}
}

// TestDispatch_RPUSH_MissType ensures RPUSH returns a type error when the key already holds a non-list.
func TestDispatch_RPUSH_MissType(t *testing.T) {
	cmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
			&resp.BulkString{B: []byte("myval")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if _, ok := out.(*resp.SimpleString); !ok {
		t.Fatalf("expected resp.SimpleString for SET command response, got %T", out)
	}

	cmd = &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
			&resp.BulkString{B: []byte("xyz")},
		},
	}

	out = Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for RPUSH command response, got %T", out)
	}
}
