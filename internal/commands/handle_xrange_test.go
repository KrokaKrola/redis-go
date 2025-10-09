package commands

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatch_Xrange_Command(t *testing.T) {
	t.Run("base scenario", func(t *testing.T) {
		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_stream")},
				&resp.BulkString{Bytes: []byte("1-0")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); ok {
			t.Fatalf("unexpected resp.Error from Dispatch response, got %T", out)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_stream")},
				&resp.BulkString{Bytes: []byte("2-0")},
				&resp.BulkString{Bytes: []byte("b")},
				&resp.BulkString{Bytes: []byte("2")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); ok {
			t.Fatalf("unexpected resp.Error from Dispatch response, got %T", out)
		}

		cmd = &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_stream")},
				&resp.BulkString{Bytes: []byte("1")},
				&resp.BulkString{Bytes: []byte("2")},
			},
		}

		out = Dispatch(cmd, store)

		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected resp.Array from Dispatch response, got %T", out)
		}

		fmt.Printf("%#v\n", arr)
	})
}
