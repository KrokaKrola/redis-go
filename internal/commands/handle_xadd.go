package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleXadd(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen < 4 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XADD command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for XADD command"}
	}

	streamEntryId, ok := cmd.ArgString(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid stream-id value for XADD command"}
	}

	restArgs := cmd.Args[2:]

	if len(restArgs)%2 != 0 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XADD command"}
	}

	fields := [][]string{}

	for i := 0; i < len(restArgs); i += 2 {
		entryKey, okKey := cmd.ArgString(2 + i)
		if !okKey {
			return &resp.Error{Msg: "ERR invalid key-pair key value for XADD command"}
		}

		entryValue, okValue := cmd.ArgString(2 + i + 1)
		if !okValue {
			return &resp.Error{Msg: "ERR invalid key-pair value for XADD command"}
		}

		fields = append(fields, []string{entryKey, entryValue})
	}

	id, ok, wrongType := store.Xadd(key, streamEntryId, fields)

	if wrongType {
		return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
	}

	if !ok {
		return &resp.Error{Msg: "ERR internal XADD command error"}
	}

	return &resp.BulkString{Bytes: []byte(id)}
}
