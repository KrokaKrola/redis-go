package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleLpop(data handlerData) resp.Value {
	argsLen := data.cmd.ArgsLen()
	if argsLen == 0 || argsLen > 2 {
		return &resp.Error{Msg: "ERR wrong number of arguments for LPOP command"}
	}

	key, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for LPOP command"}
	}

	count := 1

	if argsLen == 2 {
		count, ok = data.cmd.ArgInt(1)

		if !ok {
			return &resp.Error{Msg: "ERR invalid count value for LPOP command"}
		}

		if count <= 0 {
			return &resp.Error{Msg: "ERR invalid count value for LPOP command"}
		}
	}

	v, ok := data.store.Lpop(key, count)
	if !ok {
		return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	if v.Null || len(v.Elements) == 0 {
		if count == 1 {
			return &resp.BulkString{Null: true}
		} else {
			return &resp.Array{Null: true}
		}
	}

	if count == 1 {
		return &resp.BulkString{Bytes: []byte(v.Elements[0])}
	}

	resArray := &resp.Array{}

	for _, v := range v.Elements {
		resArray.Elements = append(resArray.Elements, &resp.BulkString{Bytes: []byte(v)})
	}

	return resArray
}
