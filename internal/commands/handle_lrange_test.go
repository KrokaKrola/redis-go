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
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("a")},
			&resp.BulkString{Bytes: []byte("b")},
			&resp.BulkString{Bytes: []byte("c")},
			&resp.BulkString{Bytes: []byte("d")},
			&resp.BulkString{Bytes: []byte("e")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store, false)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("0")},
			&resp.BulkString{Bytes: []byte("1")},
		},
	}

	out = Dispatch(cmd, store, false)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elements) != 2 {
		t.Fatalf("expected resp.Array to be of length 2, got %d", len(arr.Elements))
	}

	var got []string

	for idx, v := range arr.Elements {
		a, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
		}
		got = append(got, string(a.Bytes))
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
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("a")},
			&resp.BulkString{Bytes: []byte("b")},
			&resp.BulkString{Bytes: []byte("c")},
			&resp.BulkString{Bytes: []byte("d")},
			&resp.BulkString{Bytes: []byte("e")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store, false)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("2")},
			&resp.BulkString{Bytes: []byte("4")},
		},
	}

	out = Dispatch(cmd, store, false)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elements) != 3 {
		t.Fatalf("expected resp.Array to be of length 2, got %d", len(arr.Elements))
	}

	var got []string

	for idx, v := range arr.Elements {
		a, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
		}
		got = append(got, string(a.Bytes))
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
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("2")},
			&resp.BulkString{Bytes: []byte("4")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store, false)

	if arr := out.(*resp.Array); len(arr.Elements) != 0 && arr.Null == false {
		t.Fatalf("expected array to be Zero length and not Null, got %#v", arr)
	}
}

// TestDispatch_LRANGE_Start_Greater_Than_List_Len confirms LRANGE returns empty response if start value is greater or equal to the list len
func TestDispatch_LRANGE_Start_Greater_Than_List_Len(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("a")},
			&resp.BulkString{Bytes: []byte("b")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store, false)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("2")},
			&resp.BulkString{Bytes: []byte("2")},
		},
	}

	out = Dispatch(cmd, store, false)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elements) > 0 || arr.Null {
		t.Fatalf("expected resp.Array, to be zero length %#v", arr)
	}

	if arr.Null {
		t.Fatalf("expected resp.Array, to be not nullable %#v", arr)
	}
}

// TestDispatch_LRANGE_Start_IsGreaterThanStop confirms LRANGE returns empty response if start value is greater than stop value
func TestDispatch_LRANGE_Start_IsGreaterThanStop(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("a")},
			&resp.BulkString{Bytes: []byte("b")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store, false)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("3")},
			&resp.BulkString{Bytes: []byte("2")},
		},
	}

	out = Dispatch(cmd, store, false)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elements) > 0 {
		t.Fatalf("expected resp.Array, to be zero length %#v", arr)
	}

	if arr.Null {
		t.Fatalf("expected resp.Array, to be not nullable %#v", arr)
	}
}

// TestDispatch_LRANGE_Stop_Is_Greater_Than_List_Length confirms LRANGE returns valid list of elements
func TestDispatch_LRANGE_Stop_Is_Greater_Than_List_Length(t *testing.T) {
	cmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("a")},
			&resp.BulkString{Bytes: []byte("b")},
			&resp.BulkString{Bytes: []byte("c")},
			&resp.BulkString{Bytes: []byte("d")},
			&resp.BulkString{Bytes: []byte("e")},
		},
	}

	store := store.NewStore()

	out := Dispatch(cmd, store, false)

	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected error, got %T", out)
	}

	cmd = &Command{
		Name: LRANGE_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
			&resp.BulkString{Bytes: []byte("1")},
			&resp.BulkString{Bytes: []byte("10")},
		},
	}

	out = Dispatch(cmd, store, false)

	arr, ok := out.(*resp.Array)
	if !ok {
		t.Fatalf("expected resp.Array, got %T", out)
	}

	if len(arr.Elements) != 4 {
		t.Fatalf("expected resp.Array to be of length 4, got %d", len(arr.Elements))
	}

	var got []string

	for idx, v := range arr.Elements {
		a, ok := v.(*resp.BulkString)
		if !ok {
			t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
		}
		got = append(got, string(a.Bytes))
	}

	want := []string{"b", "c", "d", "e"}

	if !slices.Equal(want, got) {
		t.Fatalf("unexpected LRANGE response, got %#v, want %#v", got, want)
	}
}

// TestDispatch_LRANGE_NegativeRange confirms LRANGE returns valid list of elements for negative start stop values
func TestDispatch_LRANGE_NegativeRange(t *testing.T) {
	cases := []struct {
		start string
		stop  string
		want  []string
	}{
		{start: "-2", stop: "-1", want: []string{"d", "e"}},
		{start: "0", stop: "-3", want: []string{"a", "b", "c"}},
		{start: "-7", stop: "-3", want: []string{"a", "b", "c"}},
		{start: "0", stop: "-7", want: []string{}},
	}

	for _, v := range cases {
		cmd := &Command{
			Name: RPUSH_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("list_key")},
				&resp.BulkString{Bytes: []byte("a")},
				&resp.BulkString{Bytes: []byte("b")},
				&resp.BulkString{Bytes: []byte("c")},
				&resp.BulkString{Bytes: []byte("d")},
				&resp.BulkString{Bytes: []byte("e")},
			},
		}

		store := store.NewStore()

		out := Dispatch(cmd, store, false)

		if _, ok := out.(*resp.Error); ok {
			t.Fatalf("unexpected error, got %T", out)
		}

		cmd = &Command{
			Name: LRANGE_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("list_key")},
				&resp.BulkString{Bytes: []byte(v.start)},
				&resp.BulkString{Bytes: []byte(v.stop)},
			},
		}

		out = Dispatch(cmd, store, false)

		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected resp.Array, got %T", out)
		}

		if len(arr.Elements) != len(v.want) {
			t.Fatalf("expected resp.Array={%#v}, got %#v", v.want, arr.Elements)
		}

		var got []string

		for idx, v := range arr.Elements {
			a, ok := v.(*resp.BulkString)
			if !ok {
				t.Fatalf("expected idx: %d element to be *resp.BulkString, got %T", idx, a)
			}
			got = append(got, string(a.Bytes))
		}

		if !slices.Equal(v.want, got) {
			t.Fatalf("unexpected LRANGE response, got %#v, want %#v", got, v.want)
		}
	}

}
