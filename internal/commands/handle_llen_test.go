package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// TestDispatch_Llen confirms LLEN returns valid list length
func TestDispatch_Llen(t *testing.T) {
	cmd := &Command{
		Name: LPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("foo")},
			&resp.BulkString{Bytes: []byte("bar")},
			&resp.BulkString{Bytes: []byte("baz")},
		},
	}

	store := store.NewStore()

	testDispatch(cmd, store, false)

	cmd = &Command{
		Name: LLEN_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
		},
	}

	out := testDispatch(cmd, store, false)
	in, ok := out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.Number != 3 {
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.Number, 3)
	}
}

// TestDispatch_Llen_Missing_Key confirms LLEN returns 0 list length for missing key
func TestDispatch_Llen_Missing_Key(t *testing.T) {
	store := store.NewStore()

	cmd := &Command{
		Name: LLEN_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
		},
	}

	out := testDispatch(cmd, store, false)
	in, ok := out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.Number != 0 {
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.Number, 0)
	}
}

// TestDispatch_Llen_Type_Missmatch confirms resp.Error response on operation against a key holding wrong kind of value
func TestDispatch_Llen_Type_Missmatch(t *testing.T) {
	cmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("mykey")},
			&resp.BulkString{Bytes: []byte("myval")},
		},
	}

	store := store.NewStore()

	testDispatch(cmd, store, false)

	cmd = &Command{
		Name: LLEN_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("mykey")},
		},
	}

	out := testDispatch(cmd, store, false)
	_, ok := out.(*resp.Error)
	if !ok {
		t.Fatalf("expected RespError, got %T", out)
	}
}
