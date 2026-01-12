package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePush(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	argsLen := handlerCtx.Cmd.ArgsLen()
	if argsLen < 2 {
		return &resp.Error{Msg: fmt.Sprintf("ERR wrong number of arguments for %s command", handlerCtx.Cmd.Name)}
	}

	key, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: fmt.Sprintf("ERR invalid key value for %s command", handlerCtx.Cmd.Name)}
	}

	var values []string

	for argsLen-len(values) != 1 {
		value, ok := handlerCtx.Cmd.ArgString(len(values) + 1)
		if !ok {
			return &resp.Error{Msg: fmt.Sprintf("ERR invalid type of %s arguments list item", handlerCtx.Cmd.Name)}
		}

		values = append(values, value)
	}

	if len(values) == 0 {
		return &resp.Error{Msg: fmt.Sprintf("ERR empty values for %s command", handlerCtx.Cmd.Name)}
	}

	var len int64
	var isPushOk bool

	switch handlerCtx.Cmd.Name {
	case RPUSH_COMMAND:
		len, isPushOk = serverCtx.Store.Rpush(key, values)
	case LPUSH_COMMAND:
		len, isPushOk = serverCtx.Store.Lpush(key, values)
	}

	if !isPushOk {
		return &resp.Error{Msg: fmt.Sprintf("WRONGTYPE Operation against a key holding the wrong kind of value for %s command", handlerCtx.Cmd.Name)}
	}

	return &resp.Integer{Number: len}
}
