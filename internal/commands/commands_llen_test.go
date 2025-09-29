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
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("foo")},
			&resp.BulkString{B: []byte("bar")},
			&resp.BulkString{B: []byte("baz")},
		},
	}

	store := store.NewStore()

	Dispatch(cmd, store)

	cmd = &Command{
		Name: LLEN_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
		},
	}

	out := Dispatch(cmd, store)
	in, ok := out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.N != 3 {
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.N, 3)
	}
}

// TestDispatch_Llen_Missing_Key confirms LLEN returns 0 list length for missing key
func TestDispatch_Llen_Missing_Key(t *testing.T) {
	store := store.NewStore()

	cmd := &Command{
		Name: LLEN_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
		},
	}

	out := Dispatch(cmd, store)
	in, ok := out.(*resp.Integer)
	if !ok {
		t.Fatalf("expected Integer, got %T", out)
	}

	if in.N != 0 {
		t.Fatalf("unexpected LPUSH response, got %d, want %d", in.N, 0)
	}
}

// TestDispatch_Llen_Type_Missmatch confirms resp.Error response on operation against a key holding wrong kind of value
func TestDispatch_Llen_Type_Missmatch(t *testing.T) {
	cmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
			&resp.BulkString{B: []byte("myval")},
		},
	}

	store := store.NewStore()

	Dispatch(cmd, store)

	cmd = &Command{
		Name: LLEN_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("mykey")},
		},
	}

	out := Dispatch(cmd, store)
	_, ok := out.(*resp.Error)
	if !ok {
		t.Fatalf("expected RespError, got %T", out)
	}
}
