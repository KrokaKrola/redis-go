package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleBlpop(data handlerData) resp.Value {
	// TODO: add support for array-like keys
	argsLen := data.cmd.ArgsLen()

	if argsLen != 2 {
		return &resp.Error{Msg: "ERR invalid number of arguments for BLPOP command"}
	}

	key, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for BLPOP command"}
	}

	timeoutInSeconds, ok := data.cmd.ArgFloat(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid timeout value for BLPOP command"}
	}

	if timeoutInSeconds < 0 {
		return &resp.Error{Msg: "ERR invalid timeout value for BLPOP command"}
	}

	el, ok, timeout := data.store.Blpop(key, timeoutInSeconds)

	if timeout {
		return &resp.Array{Null: true}
	}

	if !ok {
		return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	return &resp.Array{Elements: []resp.Value{
		&resp.BulkString{Bytes: []byte(key)},
		&resp.BulkString{Bytes: []byte(el)},
	}}
}
