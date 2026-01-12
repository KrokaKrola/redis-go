package commands

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleInfo(data handlerData) resp.Value {
	section, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid section value"}
	}

	response := "# Server\r\n"

	if strings.EqualFold(section, "replication") {
		role := "master"

		if data.config.isReplica {
			role = "slave"
		}

		// role key-value pair
		response += "role:" + role + "r\n"

		// master_replid key-value pair
		response += fmt.Sprintf("master_replid:%s\r\n", data.replicationId)

		// master_repl_offset key-value pair
		response += fmt.Sprintf("master_repl_offset:%d", 0)

		return &resp.BulkString{
			Bytes: []byte(response),
		}
	}

	return &resp.Error{Msg: "ERR not implemented yet :("}
}
