package commands

import (
	"fmt"
	"math"
	"testing"
	"time"

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

func TestDispatch_XADD_ID_SeqNumber_Autogenerated(t *testing.T) {
	t.Run("base case positive scenario", func(t *testing.T) {
		expectedId := "1-0"
		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("1-*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		bs, ok := out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		if string(bs.Bytes) != expectedId {
			t.Fatalf("expected id=%s from Dispatch response, got %s", expectedId, bs.Bytes)
		}

		expectedId = "1-1"

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("1-*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		bs, ok = out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		if string(bs.Bytes) != expectedId {
			t.Fatalf("expected id=%s from Dispatch response, got %s", expectedId, bs.Bytes)
		}
	})

	t.Run("0 msTime value", func(t *testing.T) {
		expectedId := "0-1"
		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("0-*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		bs, ok := out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		if string(bs.Bytes) != expectedId {
			t.Fatalf("expected id=%s from Dispatch response, got %s", expectedId, bs.Bytes)
		}
	})

	t.Run("new key sequence element will get 0 seqNumber", func(t *testing.T) {
		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("1-*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		expectedId := "2-0"
		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("2-*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		bs, ok := out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		if string(bs.Bytes) != expectedId {
			t.Fatalf("expected id=%s from Dispatch response, got %s", expectedId, bs.Bytes)
		}
	})

	t.Run("seqNumber overflow error", func(t *testing.T) {
		seqId := fmt.Sprintf("%d-%d", 1, uint64(math.MaxUint64))
		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte(seqId)},
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
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("1-*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response, got %T", out)
		}
	})
}

func TestDispatch_XADD_ID_Autogenerated(t *testing.T) {
	t.Run("base case positive scenario", func(t *testing.T) {
		store := store.NewStore()

		start := uint64(time.Now().UnixMilli())

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out := Dispatch(cmd, store)

		bs, ok := out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		firstMs, firstSeq := mustParseStreamID(t, string(bs.Bytes))

		afterFirst := uint64(time.Now().UnixMilli())

		if firstMs < start || firstMs > afterFirst {
			t.Fatalf("expected auto-generated msTime to be between %d and %d, got %d", start, afterFirst, firstMs)
		}

		if firstSeq != 0 {
			t.Fatalf("expected first auto-generated seqNumber to be 0, got %d", firstSeq)
		}

		firstId := fmt.Sprintf("%d-%d", firstMs, firstSeq)

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		bs, ok = out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		secondMs, secondSeq := mustParseStreamID(t, string(bs.Bytes))

		if secondMs < firstMs {
			t.Fatalf("expected second auto-generated msTime to be >= %d from id=%s, got %d from id=%s", firstMs, firstId, secondMs, bs.Bytes)
		}

		if secondMs == firstMs {
			if secondSeq != firstSeq+1 {
				t.Fatalf("expected second auto-generated seqNumber to be %d when msTime matches, got %d", firstSeq+1, secondSeq)
			}
		} else {
			if secondSeq != 0 {
				t.Fatalf("expected seqNumber to reset to 0 when msTime advances (from %d to %d), got %d", firstMs, secondMs, secondSeq)
			}
		}
	})

	t.Run("respects future manual timestamp", func(t *testing.T) {
		store := store.NewStore()

		futureMs := uint64(time.Now().Add(10 * time.Second).UnixMilli())
		manualId := fmt.Sprintf("%d-%d", futureMs, 5)

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte(manualId)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		manualOut := Dispatch(cmd, store)

		if _, ok := manualOut.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString from manual insert, got %T", manualOut)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		autoOut := Dispatch(cmd, store)

		bs, ok := autoOut.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", autoOut)
		}

		ms, seq := mustParseStreamID(t, string(bs.Bytes))

		if ms != futureMs {
			t.Fatalf("expected msTime=%d from auto-generated id, got %d", futureMs, ms)
		}

		if seq != 6 {
			t.Fatalf("expected seqNumber=6 from auto-generated id, got %d", seq)
		}
	})

	t.Run("resets sequence when timestamp advances", func(t *testing.T) {
		store := store.NewStore()

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out := Dispatch(cmd, store)

		bs, ok := out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		firstMs, _ := mustParseStreamID(t, string(bs.Bytes))

		for time.Now().UnixMilli() <= int64(firstMs) {
			time.Sleep(time.Millisecond)
		}

		out = Dispatch(cmd, store)

		bs, ok = out.(*resp.BulkString)

		if !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response, got %T", out)
		}

		secondMs, secondSeq := mustParseStreamID(t, string(bs.Bytes))

		if secondMs <= firstMs {
			t.Fatalf("expected msTime to advance (from %d), got %d", firstMs, secondMs)
		}

		if secondSeq != 0 {
			t.Fatalf("expected seqNumber to reset to 0 when msTime advances (from %d to %d), got %d", firstMs, secondMs, secondSeq)
		}
	})

	t.Run("sequence overflow returns error", func(t *testing.T) {
		store := store.NewStore()

		futureMs := uint64(time.Now().Add(10 * time.Second).UnixMilli())
		overflowId := fmt.Sprintf("%d-%d", futureMs, uint64(math.MaxUint64))

		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte(overflowId)},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out := Dispatch(cmd, store)

		if _, ok := out.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString from Dispatch response when seeding overflow state, got %T", out)
		}

		cmd = &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("my_key")},
				&resp.BulkString{Bytes: []byte("*")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out = Dispatch(cmd, store)

		if _, ok := out.(*resp.Error); !ok {
			t.Fatalf("expected resp.Error from Dispatch response due to sequence overflow, got %T", out)
		}
	})
}

func mustParseStreamID(t *testing.T, id string) (uint64, uint64) {
	t.Helper()

	ms, seq, isAutogenSeq, isAutogen, ok := parseStreamId(id)
	if !ok {
		t.Fatalf("failed to parse stream id %q", id)
	}

	if isAutogen || isAutogenSeq {
		t.Fatalf("expected concrete stream id, got autogen flags (id=%q)", id)
	}

	return ms, seq
}
