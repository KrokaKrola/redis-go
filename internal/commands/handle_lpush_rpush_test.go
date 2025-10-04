package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// TestDispatch_LPUSH confirms LPUSH prepends elements and reports the new length.
func TestDispatch_LPUSH(t *testing.T) {
	cmd := &Command{
		Name: LPUSH_COMMAND,
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
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.N, 1)
	}

	assertListEquals(t, store, "list_key", []string{"foo"})

	cmd = &Command{
		Name: LPUSH_COMMAND,
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
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.N, 2)
	}

	assertListEquals(t, store, "list_key", []string{"bar", "foo"})
}

// TestDispatch_LPUSH_MultipleElements covers prepending multiple elements and the reported length.
func TestDispatch_LPUSH_MultipleElements(t *testing.T) {
	cmd := &Command{
		Name: LPUSH_COMMAND,
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
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.N, 3)
	}

	assertListEquals(t, store, "list_key", []string{"baz", "bar", "foo"})

	cmd = &Command{
		Name: LPUSH_COMMAND,
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
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.N, 5)
	}

	assertListEquals(t, store, "list_key", []string{"qwerty", "xyz", "baz", "bar", "foo"})
}

// TestDispatch_LPUSH_EmptyListError verifies LPUSH errors when no values are provided.
func TestDispatch_LPUSH_EmptyListError(t *testing.T) {
	cmd := &Command{
		Name: LPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)
	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for invalid amount of arguments, got %T", out)
	}

	assertListEquals(t, store, "list_key", []string{})
}

// TestDispatch_LPUSH_MissType ensures LPUSH returns a type error when the key already holds a non-list.
func TestDispatch_LPUSH_MissType(t *testing.T) {
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
		Name: LPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
			&resp.BulkString{B: []byte("xyz")},
		},
	}

	out = Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for LPUSH command response, got %T", out)
	}
}

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

	assertListEquals(t, store, "list_key", []string{"foo"})

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

	assertListEquals(t, store, "list_key", []string{"foo", "bar"})
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

	assertListEquals(t, store, "list_key", []string{"foo", "bar", "baz"})

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

	assertListEquals(t, store, "list_key", []string{"foo", "bar", "baz", "xyz", "qwerty"})
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

	assertListEquals(t, store, "list_key", nil)
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
