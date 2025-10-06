package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatchXadd_Manual_ID(t *testing.T) {
	t.Run("Sequencial pushing to stream", func(t *testing.T) {
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

		id = "1728130000003-1"

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_stream")},
				&resp.BulkString{Bytes: []byte(id)},
				&resp.BulkString{Bytes: []byte("d")},
				&resp.BulkString{Bytes: []byte("4")},
			},
		}

		out = Dispatch(cmd, store)

		bs, ok = out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		if string(bs.Bytes) != id {
			t.Fatalf("expected id=%s from Dispatch response, got %s", id, bs.Bytes)
		}

	})
	t.Run("manual ID lower ms_time than previous", func(t *testing.T) {
		id1 := "1728130000003-0"
		id2 := "1728130000002-0"
		streamKey := "my_stream"

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id1)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id2)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response, got %T", out)
		}
	})

	t.Run("Manual ID Lower Seq Number Than Previous", func(t *testing.T) {
		id1 := "1728130000003-1"
		id2 := "1728130000003-0"
		streamKey := "my_stream"

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id1)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id2)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response, got %T", out)
		}
	})

	t.Run("Manual ID SeqNumber Equal To Previous", func(t *testing.T) {
		id1 := "1728130000003-1"
		id2 := "1728130000003-0"
		streamKey := "my_stream"

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id1)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id2)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response, got %T", out)
		}
	})

	t.Run("Stream id must be greater than 0-0", func(t *testing.T) {
		id1 := "0-0"
		streamKey := "my_stream"

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id1)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response, got %T", out)
		}
	})

	t.Run("Stream id negative numbers", func(t *testing.T) {
		id1 := "-1-0"
		id2 := "0--1"
		streamKey := "my_stream"

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id1)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response with id1=%s, got %T", out, id1)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(id2)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response with id2=%s, got %T", out, id2)
		}
	})
}
