package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleType(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen != 1 {
		return &resp.Error{Msg: fmt.Sprintf("ERR wrong number of arguments for %s command", cmd.Name)}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: fmt.Sprintf("ERR invalid key value for %s command", cmd.Name)}
	}

	sv, ok := store.GetStoreRawValue(key)

	if !ok {
		return &resp.SimpleString{S: []byte("none")}
	}

	return &resp.SimpleString{S: []byte(sv.GetType())}
}
