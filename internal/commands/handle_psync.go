package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePsync(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	if serverCtx.IsReplica {
		return &resp.Error{Msg: "ERR server is replica"}
	}

	_, ok := serverCtx.ReplicasRegistry.GetReplica(handlerCtx.RemoteAddr)
	if !ok {
		return &resp.Error{Msg: "ERR replica handshake failed"}
	}

	replId := serverCtx.ReplicationId
	offset := 0

	return &resp.SimpleString{
		Bytes: fmt.Appendf(nil, "%s %s %d", "FULLRESYNC", replId, offset),
	}
}
