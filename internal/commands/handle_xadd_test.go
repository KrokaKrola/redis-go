package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatch_Xadd_Base_Scenario(t *testing.T) {
	id := "1728130000003-0"
	cmd := &Command{
		Name: XADD_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("my_stream")},
			&resp.BulkString{Bytes: []byte(id)},
			&resp.BulkString{Bytes: []byte("a")},
			&resp.BulkString{Bytes: []byte("1")},
			&resp.BulkString{Bytes: []byte("b")},
			&resp.BulkString{Bytes: []byte("2")},
			&resp.BulkString{Bytes: []byte("c")},
			&resp.BulkString{Bytes: []byte("3")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	bs, ok := out.(*resp.BulkString)

	if !ok {
		t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
	}

	if string(bs.Bytes) != id {
		t.Fatalf("expected id=%s from Dispatch response, got %s", id, bs.Bytes)
	}
}
