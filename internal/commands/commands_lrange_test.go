package commands

import (
	"slices"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// TestDispatch_LRANGE confirms LRANGE returns valid list of elements
func TestDispatch_LRANGE_Basic(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("a")},
			&resp.BulkString{B: []byte("b")},
			&resp.BulkString{B: []byte("c")},
			&resp.BulkString{B: []byte("d")},
			&resp.BulkString{B: []byte("e")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("0")},
			&resp.BulkString{B: []byte("1")},
		},
	}

	out = Dispatch(cmd, store)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elems) != 2 {
		t.Fatalf("expected resp.Array to be of length 2, got %d", len(arr.Elems))
	}

	var got []string

	for idx, v := range arr.Elems {
		a, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
		}
		got = append(got, string(a.B))
	}

	want := []string{"a", "b"}

	if !slices.Equal(want, got) {
		t.Fatalf("unexpected LRANGE response, got %#v, want %#v", got, want)
	}
}

// TestDispatch_LRANGE_Another_Basic confirms LRANGE returns valid list of elements
func TestDispatch_LRANGE_Another_Basic(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("a")},
			&resp.BulkString{B: []byte("b")},
			&resp.BulkString{B: []byte("c")},
			&resp.BulkString{B: []byte("d")},
			&resp.BulkString{B: []byte("e")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("2")},
			&resp.BulkString{B: []byte("4")},
		},
	}

	out = Dispatch(cmd, store)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elems) != 3 {
		t.Fatalf("expected resp.Array to be of length 2, got %d", len(arr.Elems))
	}

	var got []string

	for idx, v := range arr.Elems {
		a, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
		}
		got = append(got, string(a.B))
	}

	want := []string{"c", "d", "e"}

	if !slices.Equal(want, got) {
		t.Fatalf("unexpected LRANGE response, got %#v, want %#v", got, want)
	}
}

// TestDispatch_LRANGE_No_Key confirms LRANGE returns empty response if key doesn't exist
func TestDispatch_LRANGE_No_Key(t *testing.T) {
	cmd := &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("2")},
			&resp.BulkString{B: []byte("4")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if arr := out.(*resp.Array); !arr.Null || len(arr.Elems) != 0 {
		t.Fatalf("expected array to be Null length, got %#v", arr)
	}
}

// TestDispatch_LRANGE_Start_Greater_Than_List_Len confirms LRANGE returns empty response if start value is greater or equal to the list len
func TestDispatch_LRANGE_Start_Greater_Than_List_Len(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("a")},
			&resp.BulkString{B: []byte("b")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("2")},
			&resp.BulkString{B: []byte("2")},
		},
	}

	out = Dispatch(cmd, store)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elems) > 0 || !arr.Null {
		t.Fatalf("expected resp.Array, to be zero length %#v", arr)
	}
}

// TestDispatch_LRANGE_Start_IsGreaterThanStop confirms LRANGE returns empty response if start value is greater than stop value
func TestDispatch_LRANGE_Start_IsGreaterThanStop(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("a")},
			&resp.BulkString{B: []byte("b")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("3")},
			&resp.BulkString{B: []byte("2")},
		},
	}

	out = Dispatch(cmd, store)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elems) > 0 || !arr.Null {
		t.Fatalf("expected resp.Array, to be zero length %#v", arr)
	}
}

// TestDispatch_LRANGE_Stop_Is_Greater_Than_List_Length confirms LRANGE returns valid list of elements
func TestDispatch_LRANGE_Stop_Is_Greater_Than_List_Length(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("a")},
			&resp.BulkString{B: []byte("b")},
			&resp.BulkString{B: []byte("c")},
			&resp.BulkString{B: []byte("d")},
			&resp.BulkString{B: []byte("e")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("list_key")},
			&resp.BulkString{B: []byte("1")},
			&resp.BulkString{B: []byte("10")},
		},
	}

	out = Dispatch(cmd, store)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elems) != 4 {
		t.Fatalf("expected resp.Array to be of length 2, got %d", len(arr.Elems))
	}

	var got []string

	for idx, v := range arr.Elems {
		a, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
		}
		got = append(got, string(a.B))
	}

	want := []string{"b", "c", "d", "e"}

	if !slices.Equal(want, got) {
		t.Fatalf("unexpected LRANGE response, got %#v, want %#v", got, want)
	}
}
