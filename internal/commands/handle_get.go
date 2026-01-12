package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleGet(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	if handlerCtx.Cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for GET command"}
	}

	key, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for GET command"}
	}

	v, ok := serverCtx.Store.Get(key)

	if !ok {
		return &resp.BulkString{Null: true}
	}

	return &resp.BulkString{Bytes: v}
}
