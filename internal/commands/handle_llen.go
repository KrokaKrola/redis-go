package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleLlen(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	if handlerCtx.Cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for LLEN command"}
	}

	key, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for LLEN command"}
	}

	v, ok := serverCtx.Store.Lrange(key, 0, -1)
	if !ok {
		return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	if v.Null || len(v.Elements) == 0 {
		return &resp.Integer{Number: 0}
	}

	return &resp.Integer{Number: int64(len(v.Elements))}
}
