package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestHandleIncr(t *testing.T) {
	t.Run("key exists and has a numeric value", func(t *testing.T) {
		keyValue := "mykey"
		store := store.NewStore()

		cmd := &Command{
			Name: SET_COMMAND,
			Args: []resp.Value{
				&resp.SimpleString{Bytes: []byte(keyValue)},
				&resp.SimpleString{Bytes: []byte("1")},
			},
		}

		out := Dispatch(cmd, store)
		if err, ok := out.(*resp.Error); ok {
			t.Fatalf("unexpected resp.Error from SET command Dispatch response %s", err)
		}

		cmd = &Command{
			Name: INCR_COMMAND,
			Args: []resp.Value{
				&resp.SimpleString{Bytes: []byte(keyValue)},
			},
		}

		out = Dispatch(cmd, store)
		value, ok := out.(*resp.Integer)
		if !ok {
			t.Fatal("unexpected resp.Error from INCR command Dispatch response")
		}

		if value.Number != 2 {
			t.Fatalf("expected value to equal %d from Dispatch response, got %d", 2, value.Number)
		}
	})

	t.Run("key doesnt exist", func(t *testing.T) {
		keyValue := "mykey"
		store := store.NewStore()

		cmd := &Command{
			Name: INCR_COMMAND,
			Args: []resp.Value{
				&resp.SimpleString{Bytes: []byte(keyValue)},
			},
		}

		out := Dispatch(cmd, store)
		value, ok := out.(*resp.Integer)
		if !ok {
			t.Fatal("unexpected resp.Error from INCR command Dispatch response")
		}

		if value.Number != 1 {
			t.Fatalf("expected value to equal %d from Dispatch response, got %d", 1, value.Number)
		}
	})

	t.Run("key exists but doesn't have a numeric value", func(t *testing.T) {
		keyValue := "mykey"
		store := store.NewStore()

		cmd := &Command{
			Name: SET_COMMAND,
			Args: []resp.Value{
				&resp.SimpleString{Bytes: []byte(keyValue)},
				&resp.SimpleString{Bytes: []byte("not numerical value")},
			},
		}

		out := Dispatch(cmd, store)
		if err, ok := out.(*resp.Error); ok {
			t.Fatalf("unexpected resp.Error from SET command Dispatch response %s", err)
		}

		cmd = &Command{
			Name: INCR_COMMAND,
			Args: []resp.Value{
				&resp.SimpleString{Bytes: []byte(keyValue)},
			},
		}

		out = Dispatch(cmd, store)
		err, ok := out.(*resp.Error)
		if !ok {
			t.Fatal("expected resp.Error from INCR command Dispatch response")
		}

		if err.Msg != "ERR value is not an integer or out of range" {
			t.Fatal("unexpected resp.Error.Msg from INCR command Dispatch response")
		}
	})
}
