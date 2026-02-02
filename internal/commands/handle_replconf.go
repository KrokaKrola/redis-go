package commands

import (
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleReplconf(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	args := handlerCtx.Cmd.ArgsLen()

	if args <= 1 {
		return &resp.Error{Msg: "ERR invalid number of arguments for REPLCONF command"}
	}

	typeLiteral, ok := handlerCtx.Cmd.ArgString(0)

	if !ok {
		return &resp.Error{Msg: "ERR invalid command structure for REPLCONF command"}
	}

	if serverCtx.IsReplica && typeLiteral != "GETACK" {
		return &resp.Error{Msg: "ERR invalid type for REPLCONF on replica server"}
	}

	switch typeLiteral {
	case "listening-port":
		port, ok := handlerCtx.Cmd.ArgInt(1)

		if !ok {
			return &resp.Error{Msg: "ERR invalid port value"}
		}

		if err := serverCtx.ReplicasRegistry.AddReplica(handlerCtx.RemoteAddr, port); err != nil {
			return &resp.Error{Msg: err.Error()}
		}
	case "capa":
		capasList := make([]string, 0, handlerCtx.Cmd.ArgsLen()-1)
		for i := 1; i < handlerCtx.Cmd.ArgsLen(); i++ {
			capa, ok := handlerCtx.Cmd.ArgString(i)
			if ok {
				capasList = append(capasList, capa)
			} else {
				return &resp.Error{Msg: "ERR invalid capabilities value"}
			}
		}

		if err := serverCtx.ReplicasRegistry.AddCapabilities(handlerCtx.RemoteAddr, capasList); err != nil {
			return &resp.Error{Msg: err.Error()}
		}
	case "GETACK":
		if !serverCtx.IsReplica {
			return &resp.Error{Msg: "ERR server is not a replica"}
		}

		offset := 0

		return &resp.Array{
			Elements: []resp.Value{
				&resp.BulkString{Bytes: []byte(REPLCONF)},
				&resp.BulkString{Bytes: []byte("ACK")},
				&resp.BulkString{Bytes: []byte(strconv.Itoa(offset))},
			},
		}
	default:
		return &resp.Error{Msg: "ERR unknown REPLCONF command type"}
	}

	return &resp.SimpleString{Bytes: []byte("OK")}
}
