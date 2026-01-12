package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePsync(data handlerData) resp.Value {
	_, ok := data.replicasRegistry.GetReplica(data.remoteAddr)
	if !ok {
		return &resp.Error{Msg: "ERR replica handshake failed"}
	}

	replId := data.replicationId
	offset := 0

	return &resp.SimpleString{
		Bytes: fmt.Appendf(nil, "%s %s %d", "FULLRESYNC", replId, offset),
	}
}
