package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleXrange(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	argsLen := handlerCtx.Cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XRANGE command"}
	}

	key, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for XRANGE command"}
	}

	start, ok := handlerCtx.Cmd.ArgString(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid start value for XRANGE command"}
	}

	end, ok := handlerCtx.Cmd.ArgString(2)
	if !ok {
		return &resp.Error{Msg: "ERR invalid end value for XRANGE command"}
	}

	stream, err := serverCtx.Store.Xrange(key, start, end)
	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	if len(stream.Elements) == 0 {
		return &resp.Array{}
	}

	arr := populateRespArrayFromStream(stream)

	return arr
}
