package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleEcho(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	if handlerCtx.Cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for ECHO command"}
	}

	b, ok := handlerCtx.Cmd.ArgBytes(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid argument for ECHO command"}
	}

	return &resp.BulkString{Bytes: b}
}
