package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePing(serverContext *ServerContext, handlerContext *HandlerContext) resp.Value {
	argsLen := handlerContext.Cmd.ArgsLen()
	if argsLen == 0 {
		return &resp.SimpleString{Bytes: []byte("PONG")}
	}

	if argsLen == 1 {
		b, ok := handlerContext.Cmd.ArgBytes(0)
		if !ok {
			return &resp.Error{Msg: "ERR invalid argument for PING command"}
		}

		return &resp.BulkString{Bytes: b}
	}

	return &resp.Error{Msg: "ERR Invalid arguments for PING command"}
}
