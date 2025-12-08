package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handlePush(data handlerData) resp.Value {
	argsLen := data.cmd.ArgsLen()
	if argsLen < 2 {
		return &resp.Error{Msg: fmt.Sprintf("ERR wrong number of arguments for %s command", data.cmd.Name)}
	}

	key, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: fmt.Sprintf("ERR invalid key value for %s command", data.cmd.Name)}
	}

	var values []string

	for argsLen-len(values) != 1 {
		value, ok := data.cmd.ArgString(len(values) + 1)
		if !ok {
			return &resp.Error{Msg: fmt.Sprintf("ERR invalid type of %s arguments list item", data.cmd.Name)}
		}

		values = append(values, value)
	}

	if len(values) == 0 {
		return &resp.Error{Msg: fmt.Sprintf("ERR empty values for %s command", data.cmd.Name)}
	}

	var len int64
	var isPushOk bool

	switch data.cmd.Name {
	case RPUSH_COMMAND:
		len, isPushOk = data.store.Rpush(key, values)
	case LPUSH_COMMAND:
		len, isPushOk = data.store.Lpush(key, values)
	}

	if !isPushOk {
		return &resp.Error{Msg: fmt.Sprintf("WRONGTYPE Operation against a key holding the wrong kind of value for %s command", data.cmd.Name)}
	}

	return &resp.Integer{Number: len}
}
