package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleType(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	argsLen := handlerCtx.Cmd.ArgsLen()

	if argsLen != 1 {
		return &resp.Error{Msg: fmt.Sprintf("ERR wrong number of arguments for %s command", handlerCtx.Cmd.Name)}
	}

	key, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: fmt.Sprintf("ERR invalid key value for %s command", handlerCtx.Cmd.Name)}
	}

	sv, ok := serverCtx.Store.GetStoreRawValue(key)

	if !ok {
		return &resp.SimpleString{Bytes: []byte("none")}
	}

	return &resp.SimpleString{Bytes: []byte(sv.GetType())}
}
