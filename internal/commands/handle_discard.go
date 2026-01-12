package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleDiscard(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	if handlerCtx.Cmd.ArgsLen() != 0 {
		return &resp.Error{Msg: "ERR invalid number of arguments for DISCARD command"}
	}

	return &resp.SimpleString{Bytes: []byte("OK")}
}
