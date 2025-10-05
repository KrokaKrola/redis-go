package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handlePing(cmd *Command, s *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()
	if argsLen == 0 {
		return &resp.SimpleString{Bytes: []byte("PONG")}
	}

	if argsLen == 1 {
		b, ok := cmd.ArgBytes(0)
		if !ok {
			return &resp.Error{Msg: "ERR invalid argument for PING command"}
		}

		return &resp.BulkString{Bytes: b}
	}

	return &resp.Error{Msg: "ERR Invalid arguments for PING command"}
}
