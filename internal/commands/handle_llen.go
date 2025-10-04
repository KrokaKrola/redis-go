package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleLlen(cmd *Command, store *store.Store) resp.Value {
	if cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for LLEN command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for LLEN command"}
	}

	v, ok := store.Lrange(key, 0, -1)
	if !ok {
		return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	if v.Null || len(v.L) == 0 {
		return &resp.Integer{N: 0}
	}

	return &resp.Integer{N: int64(len(v.L))}
}
