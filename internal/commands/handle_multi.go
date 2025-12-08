package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleMulti(data handlerData) resp.Value {
	if data.cmd.ArgsLen() != 0 {
		return &resp.Error{Msg: "ERR invalid number of arguments for MULTI command"}
	}

	return &resp.SimpleString{Bytes: []byte("OK")}
}
