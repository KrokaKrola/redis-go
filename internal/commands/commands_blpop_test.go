package commands

import (
	"fmt"
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type blpopResult struct {
	waiterID int
	values   []string
	err      error
}

func TestDispatch_BLPopImmediateValue(t *testing.T) {
	s := store.NewStore()
	createListWithValues(t, s, "numbers", []string{"one"})

	cmd := &Command{
		Name: BLPOP_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("numbers")},
			&resp.BulkString{B: []byte("0")},
		},
	}

	out := Dispatch(cmd, s)
	arr := assertArrayResponse(t, out)

	if len(arr.Elems) != 2 {
		t.Fatalf("expected two elements in BLPOP response, got %d", len(arr.Elems))
	}

	key := asBulkString(t, arr.Elems[0])
	value := asBulkString(t, arr.Elems[1])

	if key != "numbers" {
		t.Fatalf("expected key 'numbers', got %q", key)
	}
	if value != "one" {
		t.Fatalf("expected value 'one', got %q", value)
	}

	assertListEquals(t, s, "numbers", []string{})
}

func TestDispatch_BLPopBlocksUntilPush(t *testing.T) {
	s := store.NewStore()
	resultCh := make(chan resp.Value, 1)

	go func() {
		cmd := &Command{
			Name: BLPOP_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{B: []byte("letters")},
				&resp.BulkString{B: []byte("0")},
			},
		}
		resultCh <- Dispatch(cmd, s)
	}()

	select {
	case res := <-resultCh:
		t.Fatalf("expected BLPOP to block, got response %T", res)
	case <-time.After(30 * time.Millisecond):
	}

	pushCmd := &Command{
		Name: RPUSH_COMMAND,
		Args: []resp.Value{
			&resp.BulkString{B: []byte("letters")},
			&resp.BulkString{B: []byte("x")},
		},
	}

	if out := Dispatch(pushCmd, s); isError(out) {
		t.Fatalf("unexpected error from RPUSH: %v", out)
	}

	select {
	case res := <-resultCh:
		arr := assertArrayResponse(t, res)
		if len(arr.Elems) != 2 {
			t.Fatalf("expected 2 elements, got %d", len(arr.Elems))
		}

		key := asBulkString(t, arr.Elems[0])
		value := asBulkString(t, arr.Elems[1])

		if key != "letters" {
			t.Fatalf("expected key 'letters', got %q", key)
		}
		if value != "x" {
			t.Fatalf("expected value 'x', got %q", value)
		}
	case <-time.After(time.Second):
		t.Fatalf("timed out waiting for BLPOP response after push")
	}
}

func TestDispatch_BLPopServesWaitersInOrder(t *testing.T) {
	s := store.NewStore()
	results := make(chan blpopResult, 2)
	ready := make(chan struct{}, 2)

	startWaiter := func(id int) {
		go func() {
			ready <- struct{}{}
			cmd := &Command{
				Name: BLPOP_COMMAND,
				Args: []resp.Value{
					&resp.BulkString{B: []byte("queue")},
					&resp.BulkString{B: []byte("0")},
				},
			}

			out := Dispatch(cmd, s)
			arr, err := arrayResponse(out)
			if err != nil {
				results <- blpopResult{waiterID: id, err: err}
				return
			}

			if len(arr.Elems) != 2 {
				results <- blpopResult{waiterID: id, err: fmt.Errorf("expected 2 elements, got %d", len(arr.Elems))}
				return
			}

			key := asBulkString(t, arr.Elems[0])
			value := asBulkString(t, arr.Elems[1])
			if key != "queue" {
				results <- blpopResult{waiterID: id, err: fmt.Errorf("expected key 'queue', got %q", key)}
				return
			}

			results <- blpopResult{waiterID: id, values: []string{value}}
		}()

		<-ready
		time.Sleep(10 * time.Millisecond)
	}

	startWaiter(1)
	startWaiter(2)

	push := func(value string) {
		cmd := &Command{
			Name: RPUSH_COMMAND,
			Args: []resp.Value{
				&resp.BulkString{B: []byte("queue")},
				&resp.BulkString{B: []byte(value)},
			},
		}

		if out := Dispatch(cmd, s); isError(out) {
			t.Fatalf("unexpected error from RPUSH: %v", out)
		}
	}

	push("first")

	first := <-results
	if first.err != nil {
		t.Fatalf("first waiter returned error: %v", first.err)
	}
	if first.waiterID != 1 {
		t.Fatalf("expected waiter 1 to receive first value, got waiter %d", first.waiterID)
	}
	if len(first.values) != 1 || first.values[0] != "first" {
		t.Fatalf("expected first value 'first', got %#v", first.values)
	}

	push("second")

	second := <-results
	if second.err != nil {
		t.Fatalf("second waiter returned error: %v", second.err)
	}
	if second.waiterID != 2 {
		t.Fatalf("expected waiter 2 to receive second value, got waiter %d", second.waiterID)
	}
	if len(second.values) != 1 || second.values[0] != "second" {
		t.Fatalf("expected second value 'second', got %#v", second.values)
	}
}

func assertArrayResponse(t *testing.T, v resp.Value) *resp.Array {
	t.Helper()
	arr, err := arrayResponse(v)
	if err != nil {
		t.Fatalf("%v", err)
	}
	return arr
}

func arrayResponse(v resp.Value) (*resp.Array, error) {
	arr, ok := v.(*resp.Array)
	if !ok {
		return nil, fmt.Errorf("expected *resp.Array, got %T", v)
	}
	if arr.Null {
		return nil, fmt.Errorf("unexpected null array response")
	}
	return arr, nil
}

func asBulkString(t *testing.T, v resp.Value) string {
	t.Helper()
	bs, ok := v.(*resp.BulkString)
	if !ok {
		t.Fatalf("expected *resp.BulkString, got %T", v)
	}
	if bs.Null {
		t.Fatalf("unexpected null bulk string")
	}
	return string(bs.B)
}

func isError(v resp.Value) bool {
	_, ok := v.(*resp.Error)
	return ok
}
