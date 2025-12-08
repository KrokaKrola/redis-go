package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePing(data handlerData) resp.Value {
	argsLen := data.cmd.ArgsLen()
	if argsLen == 0 {
		return &resp.SimpleString{Bytes: []byte("PONG")}
	}

	if argsLen == 1 {
		b, ok := data.cmd.ArgBytes(0)
		if !ok {
			return &resp.Error{Msg: "ERR invalid argument for PING command"}
		}

		return &resp.BulkString{Bytes: b}
	}

	return &resp.Error{Msg: "ERR Invalid arguments for PING command"}
}
