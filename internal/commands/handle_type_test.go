package commands

import (
	"testing"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func TestDispatch_Type_WrongArgCount(t *testing.T) {
	store := store.NewStore()
	cmd := &Command{Name: TYPE_COMMAND}

	out := testDispatch(cmd, store, false)
	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for missing key argument, got %T", out)
	}
}

func TestDispatch_Type_InvalidKeyType(t *testing.T) {
	store := store.NewStore()
	cmd := &Command{
		Name: TYPE_COMMAND,
		Args: []resp.Value{&resp.Integer{Number: 1}},
	}

	out := testDispatch(cmd, store, false)
	if _, ok := out.(*resp.Error); !ok {
		t.Fatalf("expected resp.Error for non-string key, got %T", out)
	}
}

func TestDispatch_Type_KeyDoesNotExist(t *testing.T) {
	store := store.NewStore()
	cmd := &Command{
		Name: TYPE_COMMAND,
		Args: []resp.Value{&resp.BulkString{Bytes: []byte("missing")}},
	}

	out := testDispatch(cmd, store, false)
	ss, ok := out.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", out)
	}

	if string(ss.Bytes) != "none" {
		t.Fatalf("unexpected TYPE response, got %q, want %q", string(ss.Bytes), "none")
	}
}

func TestDispatch_Type_StringKey(t *testing.T) {
	store := store.NewStore()
	createKeyWithValueForIndefiniteTime(t, store, "key", "value")

	typeCmd := &Command{
		Name: TYPE_COMMAND,
		Args: []resp.Value{&resp.BulkString{Bytes: []byte("key")}},
	}

	out := testDispatch(typeCmd, store, false)
	ss, ok := out.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", out)
	}

	if string(ss.Bytes) != "string" {
		t.Fatalf("unexpected TYPE response, got %q, want %q", string(ss.Bytes), "string")
	}
}

func TestDispatch_Type_ListKey(t *testing.T) {
	store := store.NewStore()
	createListWithValues(t, store, "list", []string{"a", "b"})

	typeCmd := &Command{
		Name: TYPE_COMMAND,
		Args: []resp.Value{&resp.BulkString{Bytes: []byte("list")}},
	}

	out := testDispatch(typeCmd, store, false)
	ss, ok := out.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", out)
	}

	if string(ss.Bytes) != "list" {
		t.Fatalf("unexpected TYPE response, got %q, want %q", string(ss.Bytes), "list")
	}
}

func TestDispatch_Type_ExpiredKey(t *testing.T) {
	store := store.NewStore()
	createKeyWithValueForLimitedTime(t, store, "expiring", "value", "PX", "1")

	time.Sleep(2 * time.Millisecond)

	typeCmd := &Command{
		Name: TYPE_COMMAND,
		Args: []resp.Value{&resp.BulkString{Bytes: []byte("expiring")}},
	}

	out := testDispatch(typeCmd, store, false)
	ss, ok := out.(*resp.SimpleString)
	if !ok {
		t.Fatalf("expected *resp.SimpleString, got %T", out)
	}

	if string(ss.Bytes) != "none" {
		t.Fatalf("unexpected TYPE response after expiry, got %q, want %q", string(ss.Bytes), "none")
	}
}
