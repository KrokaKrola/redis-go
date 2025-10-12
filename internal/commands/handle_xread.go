package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleXread(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD command"}
	}

	_, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR STREAMS identifier for XREAD command"}
	}

	keys := []string{}
	offset := 1

	for i := range argsLen - 1 - offset {
		key, ok := cmd.ArgString(i + offset)
		if !ok {
			return &resp.Error{Msg: "ERR key value for XREAD command"}
		}

		keys = append(keys, key)
	}

	id, ok := cmd.ArgString(argsLen - 1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid stream id value for XREAD command"}
	}

	streams, err := store.Xread(keys, id)
	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	arr := &resp.Array{}
	hasEntries := false

	for i, stream := range streams {
		if len(stream.Elements) == 0 {
			continue
		}

		hasEntries = true

		streamElements := populateRespArrayFromStream(stream)

		arr.Elements = append(arr.Elements, resp.Value(&resp.Array{
			Elements: []resp.Value{
				&resp.BulkString{Bytes: []byte(keys[i])},
				streamElements,
			},
		}))
	}

	if !hasEntries {
		return &resp.Array{Null: true}
	}

	return arr
}
