package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleInfo(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	section, ok := handlerCtx.Cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid section value"}
	}

	response := "# Server\r\n"

	if strings.EqualFold(section, "replication") {
		role := "master"

		if serverCtx.IsReplica {
			role = "slave"
		}

		// role key-value pair
		response += "role:" + role + "r\n"

		// master_replid key-value pair
		response += fmt.Sprintf("master_replid:%s\r\n", serverCtx.ReplicationId)

		// master_repl_offset key-value pair
		response += fmt.Sprintf("master_repl_offset:%d", 0)

		return &resp.BulkString{
			Bytes: []byte(response),
		}
	}

	return &resp.Error{Msg: "ERR not implemented yet :("}
}
