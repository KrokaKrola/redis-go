package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleMulti(cmd *Command, store *store.Store) resp.Value {
	if cmd.ArgsLen() != 0 {
		return &resp.Error{Msg: "ERR invalid number of arguments for MULTI command"}
	}

	return &resp.SimpleString{Bytes: []byte("OK")}
}
