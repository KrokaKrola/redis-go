package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleXread(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD command"}
	}

	cmdIdentifier, ok := cmd.ArgString(0)
	if !ok || !strings.EqualFold(cmdIdentifier, "streams") {
		return &resp.Error{Msg: "ERR invalid STREAMS identifier for XREAD command"}
	}

	if (argsLen-1)%2 != 0 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD command"}
	}

	pairsCount := (argsLen) / 2

	streamKeyIdPairs := [][]string{}

	for i := range pairsCount {
		storeKey, ok := cmd.ArgString(i + 1)
		if !ok {
			return &resp.Error{Msg: "ERR invalid stream name value for XREAD command"}
		}

		streamId, ok := cmd.ArgString(i + pairsCount + 1)
		if !ok {
			return &resp.Error{Msg: "ERR invalid stream id value for XREAD command"}
		}

		streamKeyIdPairs = append(streamKeyIdPairs, []string{storeKey, streamId})
	}

	streams, err := store.Xread(streamKeyIdPairs)
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
				&resp.BulkString{Bytes: []byte(streamKeyIdPairs[i][0])},
				streamElements,
			},
		}))
	}

	if !hasEntries {
		return &resp.Array{Null: true}
	}

	return arr
}
