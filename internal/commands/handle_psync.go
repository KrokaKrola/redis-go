package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePsync(data handlerData) resp.Value {
	if _, ok := data.replicasRegistry.GetReplica(data.remoteAddr); !ok {
		return &resp.Error{Msg: "ERR replica handshake failed"}
	}

	replId := "<REPL_ID>"
	offset := 0

	return &resp.SimpleString{
		Bytes: fmt.Appendf(nil, "%s %s %d", "FULLRESYNC", replId, offset),
	}
}
