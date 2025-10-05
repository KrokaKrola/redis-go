package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleGet(cmd *Command, store *store.Store) resp.Value {
	if cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for GET command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for GET command"}
	}

	v, ok := store.Get(key)

	if !ok {
		return &resp.BulkString{Null: true}
	}

	return &resp.BulkString{Bytes: v}
}
