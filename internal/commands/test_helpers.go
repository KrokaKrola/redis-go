package commands

import (
	"slices"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func createListWithValues(t *testing.T, s *store.Store, key string, values []string) {
	t.Helper()

	args := []resp.Value{
		&resp.BulkString{Bytes: []byte(key)},
	}

	for _, v := range values {
		args = append(args, &resp.BulkString{Bytes: []byte(v)})
	}

	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: args,
	}

	out := Dispatch(cmd, s)
	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}
}

func createKeyWithValueForIndefiniteTime(t *testing.T, store *store.Store, key string, value string) {
	setCmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte(key)},
			&resp.BulkString{Bytes: []byte(value)},
		},
	}

	if out := Dispatch(setCmd, store); true {
		ss, ok := out.(*resp.SimpleString)
		if !ok {
			t.Fatalf("expected *resp.SimpleString from SET, got %T", out)
		}
		if string(ss.Bytes) != "OK" {
			t.Fatalf("unexpected SET response, got %q, want %q", string(ss.Bytes), "OK")
		}
	}
}

func createKeyWithValueForLimitedTime(t *testing.T, store *store.Store, key string, value string, expiryType string, expiryTime string) {
	setCmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte(key)},
			&resp.BulkString{Bytes: []byte(value)},
			&resp.BulkString{Bytes: []byte(expiryType)},
			&resp.BulkString{Bytes: []byte(expiryTime)},
		},
	}

	if out := Dispatch(setCmd, store); true {
		ss, ok := out.(*resp.SimpleString)
		if !ok {
			t.Fatalf("expected *resp.SimpleString from SET, got %T", out)
		}
		if string(ss.Bytes) != "OK" {
			t.Fatalf("unexpected SET response, got %q, want %q", string(ss.Bytes), "OK")
		}
	}
}

func assertListEquals(t *testing.T, s *store.Store, key string, want []string) {
	t.Helper()

	cmd := &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte(key)},
			&resp.BulkString{Bytes: []byte("0")},
			&resp.BulkString{Bytes: []byte("-1")},
		},
	}

	out := Dispatch(cmd, s)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(want) == 0 {
		if len(arr.Elements) != 0 {
			t.Fatalf("expected zero elements for empty expectation, got %d", len(arr.Elements))
		}
		return
	}

	if arr.Null {
		t.Fatalf("expected non-null array response")
	}

	if len(arr.Elements) != len(want) {
		t.Fatalf("expected %d elements, got %d", len(want), len(arr.Elements))
	}

	got := make([]string, 0, len(arr.Elements))
	for idx, v := range arr.Elements {
		bs, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected *resp.BulkString at index %d, got %T", idx, v)
		}
		got = append(got, string(bs.Bytes))
	}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected list order, got %#v, want %#v", got, want)
	}
}

func newStreamPopulatedStore(t *testing.T, key string) *store.Store {
	t.Helper()

	s := store.NewStore()
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

	addEntriesToStore(t, s, key, entries)

	return s
}

func addEntriesToStore(t *testing.T, s *store.Store, key string, entries []struct {
	id    string
	field string
	value string
}) {
	for _, entry := range entries {
		cmd := &Command{
			Name: XADD_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte(key)},
				&resp.BulkString{Bytes: []byte(entry.id)},
				&resp.BulkString{Bytes: []byte(entry.field)},
				&resp.BulkString{Bytes: []byte(entry.value)},
			},
		}

		out := Dispatch(cmd, s)
		bs, ok := out.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected BulkString from XADD, got %T", out)
		}
		if string(bs.Bytes) != entry.id {
			t.Fatalf("expected XADD to echo id %q, got %q", entry.id, string(bs.Bytes))
		}
	}
}
