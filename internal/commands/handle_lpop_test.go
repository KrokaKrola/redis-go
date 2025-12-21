package commands

import (
	"fmt"
	"slices"
	"testing"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatch_Lpop(t *testing.T) {
	cases := []struct {
		name             string
		key              string
		listValues       []string
		count            string
		wantLpopResponse []string
		wantAfterLpop    []string
	}{
		{key: "list_key", listValues: []string{"a", "b", "c", "d"}, wantLpopResponse: []string{"a"}, wantAfterLpop: []string{"b", "c", "d"}},
		{key: "list_key", listValues: []string{"a", "b", "c", "d"}, count: "3", wantLpopResponse: []string{"a", "b", "c"}, wantAfterLpop: []string{"d"}},
		{key: "list_key", listValues: []string{"a", "b", "c", "d"}, count: "10", wantLpopResponse: []string{"a", "b", "c", "d"}, wantAfterLpop: []string{}},
	}

	for _, scenario := range cases {
		t.Run(fmt.Sprintf("lpop key: %s, values: %#v, count: %s", scenario.key, scenario.listValues, scenario.count), func(t *testing.T) {
			store := store.NewStore()
			createListWithValues(t, store, scenario.key, scenario.listValues)
			isSingleElementRequested := true

			args := []resp.Value{
				&resp.BulkString{Bytes: []byte(scenario.key)},
			}

			fmt.Printf("count: %s\n", scenario.count)

			if scenario.count != "" {
				args = append(args, &resp.BulkString{Bytes: []byte(scenario.count)})
				isSingleElementRequested = false
			}

			cmd := &Command{
				Name: LPOP_COMMAND,
				Args: args,
			}

			out := Dispatch(cmd, store, false)

			if isSingleElementRequested {
				bs, ok := out.(*resp.BulkString)
				if !ok {
					t.Fatalf("expected resp.BulkString for count=1, got %T", out)
				}

				if !slices.Equal(bs.Bytes, []byte(scenario.wantLpopResponse[0])) {
					t.Fatalf("unexpected LPOP response for count = 1, got %s, want %s", bs.Bytes, scenario.wantLpopResponse[0])
				}
			} else {

				l, ok := out.(*resp.Array)
				if !ok {
					t.Fatalf("expected resp.Array for count > 1, got %T", out)
				}

				got := []string{}

				for _, v := range l.Elements {
					bs, ok := v.(*resp.BulkString)
					if !ok {
						t.Fatalf("expected resp.BulkString, got %T", out)
					}
					got = append(got, string(bs.Bytes))
				}

				if !slices.Equal(got, scenario.wantLpopResponse) {
					t.Fatalf("unexpected LPOP response for count > 1, got %#v, want %#v", got, scenario.wantLpopResponse)
				}
			}

			assertListEquals(t, store, scenario.key, scenario.wantAfterLpop)
		})
	}
}

func TestDispatch_Lpop_EmptyList(t *testing.T) {
	store := store.NewStore()
	createListWithValues(t, store, "list_key", []string{"a"})
	cmd := &Command{
		Name: LPOP_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
		},
	}

	out := Dispatch(cmd, store, false)
	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected resp.Error, got %T", out)
	}

	out = Dispatch(cmd, store, false)
	bs, ok := out.(*resp.BulkString)
	if !ok {
		t.Fatalf("expected resp.BulkString, got %T", out)
	}

	if !bs.Null {
		t.Fatal("unexpected LPOP response for empty list, got not empty BulkString")
	}
}
func TestDispatch_Lpop_Key_Doesnt_Exist(t *testing.T) {
	store := store.NewStore()
	cmd := &Command{
		Name: LPOP_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("list_key")},
		},
	}

	out := Dispatch(cmd, store, false)
	bs, ok := out.(*resp.BulkString)
	if !ok {
		t.Fatalf("expected resp.BulkString, got %T", out)
	}

	if !bs.Null {
		t.Fatal("unexpected LPOP response for empty list, got not empty BulkString")
	}
}

func TestDispatch_Lpop_Key_Wrong_Type(t *testing.T) {
	store := store.NewStore()

	cmd := &Command{
		Name: SET_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("key")},
			&resp.BulkString{Bytes: []byte("value")},
		},
	}

	out := Dispatch(cmd, store, false)
	if _, ok := out.(*resp.Error); ok {
		t.Fatalf("unexpected resp.Error, got %T", out)
	}

	cmd = &Command{
		Name: LPOP_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{Bytes: []byte("key")},
		},
	}

	out = Dispatch(cmd, store, false)
	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error, got %T", out)
	}
}

func TestDispatch_Lpop_InvalidCount(t *testing.T) {
	cases := []struct {
		name  string
		count string
	}{
		{name: "zero", count: "0"},
		{name: "negative", count: "-1"},
	}

	for _, scenario := range cases {
		t.Run(scenario.name, func(t *testing.T) {
			store := store.NewStore()
			createListWithValues(t, store, "list_key", []string{"a"})

			cmd := &Command{
				Name: LPOP_COMMAND,
				Args: []resp.Value{
					&resp.BulkString{Bytes: []byte("list_key")},
					&resp.BulkString{Bytes: []byte(scenario.count)},
				},
			}

			out := Dispatch(cmd, store, false)

			respErr, ok := out.(*resp.Error)
			if !ok {
				t.Fatalf("expected resp.Error, got %T", out)
			}

			if respErr.Msg != "ERR invalid count value for LPOP command" {
				t.Fatalf("unexpected error message, got %q", respErr.Msg)
			}
		})
	}
}

func TestDispatch_Lpop_CountGreaterThanOneEmptyResult(t *testing.T) {
	t.Run("existing key becomes empty", func(t *testing.T) {
		store := store.NewStore()
		createListWithValues(t, store, "list_key", []string{"a"})

		cmd := &Command{
			Name: LPOP_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("list_key")},
			},
		}

		out := Dispatch(cmd, store, false)
		if _, ok := out.(*resp.BulkString); !ok {
			t.Fatalf("expected resp.BulkString when popping existing value, got %T", out)
		}

		cmd = &Command{
			Name: LPOP_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("list_key")},
				&resp.BulkString{Bytes: []byte("2")},
			},
		}

		out = Dispatch(cmd, store, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected resp.Array for count > 1, got %T", out)
		}

		if !arr.Null {
			t.Fatalf("expected null array for empty list")
		}

		if len(arr.Elements) != 0 {
			t.Fatalf("expected zero elements for null array response, got %d", len(arr.Elements))
		}
	})

	t.Run("missing key", func(t *testing.T) {
		store := store.NewStore()
		cmd := &Command{
			Name: LPOP_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{Bytes: []byte("missing")},
				&resp.BulkString{Bytes: []byte("2")},
			},
		}

		out := Dispatch(cmd, store, false)
		arr, ok := out.(*resp.Array)
		if !ok {
			t.Fatalf("expected resp.Array for missing key, got %T", out)
		}

		if !arr.Null {
			t.Fatalf("expected null array for missing key")
		}

		if len(arr.Elements) != 0 {
			t.Fatalf("expected zero elements for null array response, got %d", len(arr.Elements))
		}
	})
}
