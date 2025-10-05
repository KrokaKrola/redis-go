package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleLrange(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR wrong number of arguments for LRANGE command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for LRANGE command"}
	}

	start, ok := cmd.ArgInt(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid start value for LRANGE command"}
	}

	stop, ok := cmd.ArgInt(2)
	if !ok {
		return &resp.Error{Msg: "ERR invalid stop value for LRANGE command"}
	}

	v, ok := store.Lrange(key, start, stop)
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
