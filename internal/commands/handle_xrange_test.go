package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatchXrangeCommand(t *testing.T) {
	const streamKey = "my_stream"

	t.Run("filters elements within equal millisecond boundaries", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1-1")},
				&resp.BulkString{Bytes: []byte("1-2")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XRANGE, got %T", out)
		}

		if len(arr.Elements) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
		}

		firstEntry, ok := arr.Elements[0].(*resp.Array)
		if !ok || len(firstEntry.Elements) != 2 {
			t.Fatalf("unexpected first entry layout: %#v", arr.Elements[0])
		}
		id, ok := firstEntry.Elements[0].(*resp.BulkString)
		if !ok || string(id.Bytes) != "1-1" {
			t.Fatalf("expected first id 1-1, got %#v", firstEntry.Elements[0])
		}
		fields, ok := firstEntry.Elements[1].(*resp.Array)
		if !ok || len(fields.Elements) != 2 {
			t.Fatalf("expected field array of length 2, got %#v", firstEntry.Elements[1])
		}
		name, ok := fields.Elements[0].(*resp.BulkString)
		if !ok || string(name.Bytes) != "b" {
			t.Fatalf("expected first field name b, got %#v", fields.Elements[0])
		}
		val, ok := fields.Elements[1].(*resp.BulkString)
		if !ok || string(val.Bytes) != "1" {
			t.Fatalf("expected first field value 1, got %#v", fields.Elements[1])
		}

		secondEntry, ok := arr.Elements[1].(*resp.Array)
		if !ok || len(secondEntry.Elements) != 2 {
			t.Fatalf("unexpected second entry layout: %#v", arr.Elements[1])
		}
		id, ok = secondEntry.Elements[0].(*resp.BulkString)
		if !ok || string(id.Bytes) != "1-2" {
			t.Fatalf("expected second id 1-2, got %#v", secondEntry.Elements[0])
		}
	})

	t.Run("returns only matching element when start equals end", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1-2")},
				&resp.BulkString{Bytes: []byte("1-2")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 1 {
			t.Fatalf("expected 1 element, got %d", len(arr.Elements))
		}
		entry, ok := arr.Elements[0].(*resp.Array)
		if !ok || len(entry.Elements) != 2 {
			t.Fatalf("unexpected XRANGE entry payload: %#v", arr.Elements[0])
		}
		id, ok := entry.Elements[0].(*resp.BulkString)
		if !ok || string(id.Bytes) != "1-2" {
			t.Fatalf("expected entry id 1-2, got %#v", entry.Elements[0])
		}
	})

	t.Run("expands auto sequence bounds when timestamps equal", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1")},
				&resp.BulkString{Bytes: []byte("1")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 4 {
			t.Fatalf("expected 4 elements, got %d", len(arr.Elements))
		}
		want := []string{"1-0", "1-1", "1-2", "1-3"}
		for i, expected := range want {
			entry, ok := arr.Elements[i].(*resp.Array)
			if !ok || len(entry.Elements) != 2 {
				t.Fatalf("unexpected XRANGE entry at %d: %#v", i, arr.Elements[i])
			}
			id, ok := entry.Elements[0].(*resp.BulkString)
			if !ok || string(id.Bytes) != expected {
				t.Fatalf("entry %d: expected id %s, got %#v", i, expected, entry.Elements[0])
			}
		}
	})

	t.Run("uses auto sequence lower bound when no start sequence provided", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1")},
				&resp.BulkString{Bytes: []byte("1-1")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
		}
		firstEntry, ok := arr.Elements[0].(*resp.Array)
		if !ok || len(firstEntry.Elements) != 2 {
			t.Fatalf("unexpected first entry layout: %#v", arr.Elements[0])
		}
		secondEntry, ok := arr.Elements[1].(*resp.Array)
		if !ok || len(secondEntry.Elements) != 2 {
			t.Fatalf("unexpected second entry layout: %#v", arr.Elements[1])
		}
		id0, ok := firstEntry.Elements[0].(*resp.BulkString)
		if !ok {
			t.Fatalf("expected BulkString id, got %#v", firstEntry.Elements[0])
		}
		id1, ok := secondEntry.Elements[0].(*resp.BulkString)
		if !ok {
			t.Fatalf("expected BulkString id, got %#v", secondEntry.Elements[0])
		}
		if string(id0.Bytes) != "1-0" || string(id1.Bytes) != "1-1" {
			t.Fatalf("expected ids [1-0 1-1], got [%s %s]", id0.Bytes, id1.Bytes)
		}
	})

	t.Run("returns intermediate milliseconds inclusively", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1-2")},
				&resp.BulkString{Bytes: []byte("2-0")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 3 {
			t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
		}
		want := []string{"1-2", "1-3", "2-0"}
		for i, expected := range want {
			entry, ok := arr.Elements[i].(*resp.Array)
			if !ok || len(entry.Elements) != 2 {
				t.Fatalf("unexpected XRANGE entry at %d: %#v", i, arr.Elements[i])
			}
			id, ok := entry.Elements[0].(*resp.BulkString)
			if !ok || string(id.Bytes) != expected {
				t.Fatalf("entry %d: expected id %s, got %#v", i, expected, entry.Elements[0])
			}
		}
	})

	t.Run("returns empty array when stream missing", func(t *testing.T) {
		s := store.NewStore()

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("missing")},
				&resp.BulkString{Bytes: []byte("0-0")},
				&resp.BulkString{Bytes: []byte("9-9")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 0 {
			t.Fatalf("expected empty result, got %d", len(arr.Elements))
		}
	})

	t.Run("returns error when key holds non stream value", func(t *testing.T) {
		s := store.NewStore()

		setCmd := &Command{
			Name: SET_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("plain_key")},
				&resp.BulkString{Bytes: []byte("value")},
			},
		}
		out := testDispatch(setCmd, s, false)
		ss, ok := out.(*resp.SimpleString)
		if !ok || string(ss.Bytes) != "OK" {
			t.Fatalf("expected SET to return +OK, got %#v", out)
		}

		xrangeCmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("plain_key")},
				&resp.BulkString{Bytes: []byte("0-0")},
				&resp.BulkString{Bytes: []byte("9-9")},
			},
		}

		result := testDispatch(xrangeCmd, s, false)
		errResp, ok := result.(*resp.Error)
		if !ok {
			t.Fatalf("expected Error, got %T", result)
		}
		if errResp.Msg != "MISSTYPE of the element in the underlying stream" {
			t.Fatalf("unexpected error message: %q", errResp.Msg)
		}
	})

	t.Run("returns error for invalid start stream id", func(t *testing.T) {
		s := store.NewStore()

		addCmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("stream")},
				&resp.BulkString{Bytes: []byte("1-0")},
				&resp.BulkString{Bytes: []byte("field")},
				&resp.BulkString{Bytes: []byte("value")},
			},
		}
		out := testDispatch(addCmd, s, false)
		seed, ok := out.(*resp.BulkString)
		if !ok || string(seed.Bytes) != "1-0" {
			t.Fatalf("expected XADD to echo id 1-0, got %#v", out)
		}

		xrangeCmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("stream")},
				&resp.BulkString{Bytes: []byte("bad-start")},
				&resp.BulkString{Bytes: []byte("0-0")},
			},
		}

		result := testDispatch(xrangeCmd, s, false)
		errResp, ok := result.(*resp.Error)
		if !ok {
			t.Fatalf("expected Error, got %T", result)
		}
		if errResp.Msg != "ERR invalid start value for XRANGE command" {
			t.Fatalf("unexpected error message: %q", errResp.Msg)
		}
	})

	t.Run("returns range of stream with - as start bound", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("-")},
				&resp.BulkString{Bytes: []byte("1-3")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 4 {
			t.Fatalf("expected 4 elements, got %d", len(arr.Elements))
		}
		want := []string{"1-0", "1-1", "1-2", "1-3"}
		for i, expected := range want {
			entry, ok := arr.Elements[i].(*resp.Array)
			if !ok {
				t.Fatalf("unexpected XRANGE entry at %d: %#v", i, arr.Elements[i])
			}
			id, ok := entry.Elements[0].(*resp.BulkString)
			if !ok || string(id.Bytes) != expected {
				t.Fatalf("entry %d: expected id %s, got %#v", i, expected, entry.Elements[0])
			}
		}
	})

	t.Run("returns range of stream with + as end bound", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1-1")},
				&resp.BulkString{Bytes: []byte("+")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 4 {
			t.Fatalf("expected 4 elements, got %d", len(arr.Elements))
		}
		want := []string{"1-1", "1-2", "1-3", "2-0"}
		for i, expected := range want {
			entry, ok := arr.Elements[i].(*resp.Array)
			if !ok {
				t.Fatalf("unexpected XRANGE entry at %d: %#v", i, arr.Elements[i])
			}
			id, ok := entry.Elements[0].(*resp.BulkString)
			if !ok || string(id.Bytes) != expected {
				t.Fatalf("entry %d: expected id %s, got %#v", i, expected, entry.Elements[0])
			}
		}
	})

	t.Run("returns full range of stream with - as start bound and + as end bound", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("-")},
				&resp.BulkString{Bytes: []byte("+")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 5 {
			t.Fatalf("expected 5 elements, got %d", len(arr.Elements))
		}
		want := []string{"1-0", "1-1", "1-2", "1-3", "2-0"}
		for i, expected := range want {
			entry, ok := arr.Elements[i].(*resp.Array)
			if !ok {
				t.Fatalf("unexpected XRANGE entry at %d: %#v", i, arr.Elements[i])
			}
			id, ok := entry.Elements[0].(*resp.BulkString)
			if !ok || string(id.Bytes) != expected {
				t.Fatalf("entry %d: expected id %s, got %#v", i, expected, entry.Elements[0])
			}
		}
	})

	t.Run("returns empty stream for + +", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("+")},
				&resp.BulkString{Bytes: []byte("+")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 0 {
			t.Fatalf("expected 0 elements, got %d", len(arr.Elements))
		}
	})

	t.Run("returns empty stream for + -", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("+")},
				&resp.BulkString{Bytes: []byte("-")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 0 {
			t.Fatalf("expected 0 elements, got %d", len(arr.Elements))
		}
	})

	t.Run("returns empty stream for - -", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("-")},
				&resp.BulkString{Bytes: []byte("-")},
			},
		}

		out := testDispatch(cmd, s, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array, got %T", out)
		}
		if len(arr.Elements) != 0 {
			t.Fatalf("expected 0 elements, got %d", len(arr.Elements))
		}
	})
}
