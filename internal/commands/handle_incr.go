package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleIncr(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen != 1 {
		return &resp.Error{Msg: "ERR invalid number of arguments for INCR command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for INCR command"}
	}

	value, err := store.Incr(key)

	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	return &resp.Integer{Number: value}
}
