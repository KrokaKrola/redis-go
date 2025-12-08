package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleLrange(data handlerData) resp.Value {
	argsLen := data.cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR wrong number of arguments for LRANGE command"}
	}

	key, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for LRANGE command"}
	}

	start, ok := data.cmd.ArgInt(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid start value for LRANGE command"}
	}

	stop, ok := data.cmd.ArgInt(2)
	if !ok {
		return &resp.Error{Msg: "ERR invalid stop value for LRANGE command"}
	}

	v, ok := data.store.Lrange(key, start, stop)
	if !ok {
		return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	if v.Null {
		return &resp.Array{Null: true}
	}

	if len(v.Elements) == 0 {
		return &resp.Array{}
	}

	resArray := &resp.Array{}

	for _, v := range v.Elements {
		resArray.Elements = append(resArray.Elements, &resp.BulkString{Bytes: []byte(v)})
	}

	return resArray
}
