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
		&resp.BulkString{B: []byte(key)},
	}

	for _, v := range values {
		args = append(args, &resp.BulkString{B: []byte(v)})
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

func assertListEquals(t *testing.T, s *store.Store, key string, want []string) {
	t.Helper()

	cmd := &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte(key)},
			&resp.BulkString{B: []byte("0")},
			&resp.BulkString{B: []byte("-1")},
		},
	}

	out := Dispatch(cmd, s)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(want) == 0 {
		if !arr.Null {
			t.Fatalf("expected null array for empty expectation, got %#v", arr)
		}
		if len(arr.Elems) != 0 {
			t.Fatalf("expected zero elements for empty expectation, got %d", len(arr.Elems))
		}
		return
	}

	if arr.Null {
		t.Fatalf("expected non-null array response")
	}

	if len(arr.Elems) != len(want) {
		t.Fatalf("expected %d elements, got %d", len(want), len(arr.Elems))
	}

	got := make([]string, 0, len(arr.Elems))
	for idx, v := range arr.Elems {
		bs, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected *resp.BulkString at index %d, got %T", idx, v)
		}
		got = append(got, string(bs.B))
	}

	if !slices.Equal(got, want) {
		t.Fatalf("unexpected list order, got %#v, want %#v", got, want)
	}
}
