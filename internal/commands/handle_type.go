package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleType(data handlerData) resp.Value {
	argsLen := data.cmd.ArgsLen()

	if argsLen != 1 {
		return &resp.Error{Msg: fmt.Sprintf("ERR wrong number of arguments for %s command", data.cmd.Name)}
	}

	key, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: fmt.Sprintf("ERR invalid key value for %s command", data.cmd.Name)}
	}

	sv, ok := data.store.GetStoreRawValue(key)

	if !ok {
		return &resp.SimpleString{Bytes: []byte("none")}
	}

	return &resp.SimpleString{Bytes: []byte(sv.GetType())}
}
