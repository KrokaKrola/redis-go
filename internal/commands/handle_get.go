package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleGet(data handlerData) resp.Value {
	if data.cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for GET command"}
	}

	key, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for GET command"}
	}

	v, ok := data.store.Get(key)

	if !ok {
		return &resp.BulkString{Null: true}
	}

	return &resp.BulkString{Bytes: v}
}
