package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleInfo(data handlerData) resp.Value {
	section, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid section value"}
	}

	if strings.EqualFold(section, "replication") {
		role := "master"

		if data.config.isReplica {
			role = "slave"
		}

		return &resp.BulkString{
			Bytes: []byte("role:" + role),
		}
	}

	return &resp.Error{Msg: "ERR not implemented yet :("}
}
