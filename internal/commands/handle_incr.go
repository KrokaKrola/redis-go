package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleIncr(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	argsLen := handlerCtx.Cmd.ArgsLen()

	if argsLen != 1 {
		return &resp.Error{Msg: "ERR invalid number of arguments for INCR command"}
	}

	key, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for INCR command"}
	}

	value, err := serverCtx.Store.Incr(key)

	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	return &resp.Integer{Number: value}
}
