package commands

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

func handleReplconf(data handlerData) resp.Value {
	if data.config.isReplica {
		return &resp.Error{Msg: "ERR attempt to REPLCONF on replica server"}
	}

	args := data.cmd.ArgsLen()

	if args <= 1 {
		return &resp.Error{Msg: "ERR invalid number of arguments for REPLCONF command"}
	}

	typeLiteral, ok := data.cmd.ArgString(0)

	if !ok {
		return &resp.Error{Msg: "ERR invalid command structure for REPLCONF command"}
	}

	switch typeLiteral {
	case "listening-port":
		port, ok := data.cmd.ArgInt(1)

		if !ok {
			return &resp.Error{Msg: "ERR invalid port value"}
		}

		if err := data.replicasRegistry.AddReplica(data.remoteAddr, port); err != nil {
			return &resp.Error{Msg: err.Error()}
		}
	case "capa":
		capasList := []string{}
		for i := 1; i < data.cmd.ArgsLen(); i++ {
			capa, ok := data.cmd.ArgString(i)
			if ok {
				capasList = append(capasList, capa)
			} else {
				return &resp.Error{Msg: "ERR invalid capabilities value"}
			}
		}

		if err := data.replicasRegistry.AddCapabilities(data.remoteAddr, capasList); err != nil {
			return &resp.Error{Msg: err.Error()}
		}
	}

	return &resp.SimpleString{Bytes: []byte("OK")}
}
