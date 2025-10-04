package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handlePush(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()
	if argsLen < 2 {
		return &resp.Error{Msg: fmt.Sprintf("ERR wrong number of arguments for %s command", cmd.Name)}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: fmt.Sprintf("ERR invalid key value for %s command", cmd.Name)}
	}

	var values []string

	for argsLen-len(values) != 1 {
		value, ok := cmd.ArgString(len(values) + 1)
		if !ok {
			return &resp.Error{Msg: fmt.Sprintf("ERR invalid type of %s arguments list item", cmd.Name)}
		}

		values = append(values, value)
	}

	if len(values) == 0 {
		return &resp.Error{Msg: fmt.Sprintf("ERR empty values for %s command", cmd.Name)}
	}

	var len int64
	var isPushOk bool

	switch cmd.Name {
	case RPUSH_COMMAND:
		len, isPushOk = store.Rpush(key, values)
	case LPUSH_COMMAND:
		len, isPushOk = store.Lpush(key, values)
	}

	if !isPushOk {
		return &resp.Error{Msg: fmt.Sprintf("WRONGTYPE Operation against a key holding the wrong kind of value for %s command", cmd.Name)}
	}

	return &resp.Integer{N: len}
}
