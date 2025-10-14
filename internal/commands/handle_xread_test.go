package commands

import (
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatchXreadCommand(t *testing.T) {
	const streamKey = "my_stream"
	const streamKey2 = "my_stream_2"

	validateEntry := func(t *testing.T, value resp.Value, key string, elements []struct {
		id     string
		fields []string
	}) {
		t.Helper()
		entry, ok := value.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XREAD, got %T", entry)
		}

		if len(entry.Elements) != 2 {
			t.Fatalf("expected entry to have 2 elements, got %d", len(entry.Elements))
		}

		entryId, ok := entry.Elements[0].(*resp.BulkString)
		if !ok || string(entryId.Bytes) != key {
			t.Fatalf("expected entry[0] to equal %s, got %s", key, entryId.Bytes)
		}

		entryList, ok := entry.Elements[1].(*resp.Array)
		if !ok {
			t.Fatalf("expected entry[1] to be array, got %T", entryList)
		}

		if len(entryList.Elements) != len(elements) {
			t.Fatalf("expected entryList to have %d elements, got %d", len(elements), len(entryList.Elements))
		}

		for elIdx, el := range entryList.Elements {
			streamData, ok := el.(*resp.Array)
			if !ok {
				t.Fatalf("expected streamData to be of type Array, got %T", streamData)
			}

			streamDataId, ok := streamData.Elements[0].(*resp.BulkString)
			if !ok || string(streamDataId.Bytes) != elements[elIdx].id {
				t.Fatalf("expected streamData[%d] to equal %s, got %s", elIdx, elements[elIdx].id, streamDataId.Bytes)
			}

			streamFields, ok := streamData.Elements[1].(*resp.Array)
			if !ok {
				t.Fatalf("expected streamData.fields to be of type Array, got %T", streamFields)
			}

			if len(streamFields.Elements) != len(elements[elIdx].fields) {
				t.Fatalf("expected streamFields.Elements to have %d elements, got %d", len(elements[elIdx].fields), len(streamFields.Elements))
			}

			for fieldIdx, field := range streamFields.Elements {
				fieldAsStr, ok := field.(*resp.BulkString)
				if !ok || string(fieldAsStr.Bytes) != elements[elIdx].fields[fieldIdx] {
					t.Fatalf("expected field[%d] to equal %s, got %s", fieldIdx, elements[elIdx].fields[fieldIdx], fieldAsStr.Bytes)
				}
			}
		}
	}

	t.Run("get elements exclusively from single stream", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1-2")},
			},
		}

		out := Dispatch(cmd, s)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XREAD, got %T", out)
		}

		if len(arr.Elements) != 1 {
			t.Fatalf("expected 1 element, got %d", len(arr.Elements))
		}

		entry, ok := arr.Elements[0].(*resp.Array)
		if !ok || len(entry.Elements) != 2 {
			t.Fatalf("expected entry to have 2 elements, got %d", len(entry.Elements))
		}

		validateEntry(t, entry, streamKey, []struct {
			id     string
			fields []string
		}{
			{id: "1-3", fields: []string{"d", "3"}},
			{id: "2-0", fields: []string{"e", "4"}},
		})
	})

	t.Run("returns full range of elements from single stream", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("0-0")},
			},
		}

		out := Dispatch(cmd, s)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected resp.Array from XREAD, got %T", out)
		}

		entry, ok := arr.Elements[0].(*resp.Array)
		if !ok {
			t.Fatalf("expected resp.Array from arr.Elements[0], got %T", out)
		}

		validateEntry(t, entry, streamKey, []struct {
			id     string
			fields []string
		}{
			{id: "1-0", fields: []string{"a", "0"}},
			{id: "1-1", fields: []string{"b", "1"}},
			{id: "1-2", fields: []string{"c", "2"}},
			{id: "1-3", fields: []string{"d", "3"}},
			{id: "2-0", fields: []string{"e", "4"}},
		})
	})

	t.Run("returns null array when stream has no elements", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("2-0")},
			},
		}

		out := Dispatch(cmd, s)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XREAD, got %T", out)
		}

		if !arr.Null {
			t.Fatalf("expected Array from XREAD to be null, got %#v", out)
		}
	})

	t.Run("returns null array when stream is missing", func(t *testing.T) {
		s := store.NewStore()

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("1-2")},
			},
		}

		out := Dispatch(cmd, s)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XREAD, got %T", out)
		}

		if !arr.Null {
			t.Fatalf("expected Array from XREAD to be null, got %#v", out)
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
		out := Dispatch(setCmd, s)
		ss, ok := out.(*resp.SimpleString)
		if !ok || string(ss.Bytes) != "OK" {
			t.Fatalf("expected SET to return +OK, got %#v", out)
		}

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte("plain_key")},
				&resp.BulkString{Bytes: []byte("1-2")},
			},
		}

		result := Dispatch(cmd, s)
		errResp, ok := result.(*resp.Error)
		if !ok {
			t.Fatalf("expected Error, got %T", result)
		}
		if errResp.Msg != "MISSTYPE of the element in the underlying stream" {
			t.Fatalf("unexpected error message: %q", errResp.Msg)
		}
	})

	t.Run("returns error when user passes invalid number of arguments", func(t *testing.T) {
		s := store.NewStore()

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte("key_1")},
				&resp.BulkString{Bytes: []byte("1-2")},
				&resp.BulkString{Bytes: []byte("key_2")},
			},
		}

		result := Dispatch(cmd, s)

		errResp, ok := result.(*resp.Error)
		if !ok {
			t.Fatalf("expected Error, got %T", result)
		}
		if errResp.Msg != "ERR invalid number of arguments for XREAD command" {
			t.Fatalf("unexpected error message: %q", errResp.Msg)
		}
	})

	t.Run("returns valid entries for multiple streams key-id values", func(t *testing.T) {
		s := newStreamPopulatedStore(t, streamKey)

		entries := []struct {
			id    string
			field string
			value string
		}{
			{"1-0", "a", "0"},
			{"1-1", "b", "1"},
			{"1-2", "c", "2"},
			{"1-3", "d", "3"},
			{"2-0", "e", "4"},
		}

		addEntriesToStore(t, s, streamKey2, entries)

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte(streamKey2)},
				&resp.BulkString{Bytes: []byte("1-2")},
				&resp.BulkString{Bytes: []byte("1-3")},
			},
		}

		out := Dispatch(cmd, s)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XREAD, got %T", out)
		}

		if len(arr.Elements) != 2 {
			t.Fatalf("expected 2 element, got %d", len(arr.Elements))
		}

		entry1, ok := arr.Elements[0].(*resp.Array)
		if !ok || len(entry1.Elements) != 2 {
			t.Fatalf("expected entry to have 2 elements, got %d", len(entry1.Elements))
		}

		validateEntry(t, entry1, streamKey, []struct {
			id     string
			fields []string
		}{
			{id: "1-3", fields: []string{"d", "3"}},
			{id: "2-0", fields: []string{"e", "4"}},
		})

		entry2, ok := arr.Elements[1].(*resp.Array)
		if !ok || len(entry2.Elements) != 2 {
			t.Fatalf("expected entry to have 2 element, got %d", len(entry2.Elements))
		}

		validateEntry(t, entry2, streamKey2, []struct {
			id     string
			fields []string
		}{
			{id: "2-0", fields: []string{"e", "4"}},
		})
	})

	t.Run("returns entries for existing stream when preceding stream missing", func(t *testing.T) {
		const missingStreamKey = "missing_stream"

		s := newStreamPopulatedStore(t, streamKey)

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte(missingStreamKey)},
				&resp.BulkString{Bytes: []byte(streamKey)},
				&resp.BulkString{Bytes: []byte("0-0")},
				&resp.BulkString{Bytes: []byte("0-0")},
			},
		}

		out := Dispatch(cmd, s)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected Array from XREAD, got %T", out)
		}

		if arr.Null {
			t.Fatalf("expected Array from XREAD to be non-null")
		}

		if len(arr.Elements) != 1 {
			t.Fatalf("expected 1 entry from XREAD, got %d", len(arr.Elements))
		}

		validateEntry(t, arr.Elements[0], streamKey, []struct {
			id     string
			fields []string
		}{
			{id: "1-0", fields: []string{"a", "0"}},
			{id: "1-1", fields: []string{"b", "1"}},
			{id: "1-2", fields: []string{"c", "2"}},
			{id: "1-3", fields: []string{"d", "3"}},
			{id: "2-0", fields: []string{"e", "4"}},
		})
	})

	t.Run("blocking read is correctly parsed", func(t *testing.T) {
		s := store.NewStore()

		cmd := &Command{
			Name: XREAD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("BLOCK")},
				&resp.BulkString{Bytes: []byte("100")},
				&resp.BulkString{Bytes: []byte("STREAMS")},
				&resp.BulkString{Bytes: []byte("key")},
				&resp.BulkString{Bytes: []byte("0-0")},
			},
		}

		out := Dispatch(cmd, s)

		if arr, ok := out.(*resp.Error); ok {
			t.Fatalf("unexpected error from XREAD response %#v, expected resp.Array", arr)
		}
	})
}
