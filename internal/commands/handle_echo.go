package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleEcho(cmd *Command, s *store.Store) resp.Value {
	if cmd.ArgsLen() != 1 {
		return &resp.Error{Msg: "ERR wrong number of arguments for ECHO command"}
	}

	b, ok := cmd.ArgBytes(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid argument for ECHO command"}
	}

	return &resp.BulkString{Bytes: b}
}
