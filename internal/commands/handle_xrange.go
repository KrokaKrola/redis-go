package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleXrange(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XRANGE command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for XRANGE command"}
	}

	start, ok := cmd.ArgString(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid start value for XRANGE command"}
	}

	end, ok := cmd.ArgString(2)
	if !ok {
		return &resp.Error{Msg: "ERR invalid end value for XRANGE command"}
	}

	stream, err := store.Xrange(key, start, end)
	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	if len(stream.Elements) == 0 {
		return &resp.Array{}
	}

	arr := populateRespArrayFromStream(stream)

	return arr
}
