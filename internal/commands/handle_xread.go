package commands

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleXread(data handlerData) resp.Value {
	argsLen := data.cmd.ArgsLen()

	if argsLen < 1 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD command"}
	}

	cmdIdentifier, ok := data.cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid identifier for XREAD command"}
	}

	isBlocking := strings.EqualFold(cmdIdentifier, "block")

	if isBlocking && argsLen < 5 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD BLOCK command"}
	} else if !isBlocking && argsLen < 3 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD STREAMS command"}
	}

	blockingTimeoutMs := 0

	if isBlocking {
		blockingTimeoutMs, ok = data.cmd.ArgInt(1)

		if !ok {
			return &resp.Error{Msg: "ERR invalid BLOCK timeout value for XREAD command"}
		}

		streamsKeyword, ok := data.cmd.ArgString(2)
		if !ok || !strings.EqualFold(streamsKeyword, "streams") {
			return &resp.Error{Msg: "ERR invalid STREAMS identifier for XREAD command"}
		}
	}

	if (argsLen-1)%2 != 0 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XREAD command"}
	}

	pairsCount := (argsLen) / 2

	if isBlocking {
		pairsCount = pairsCount - 1
	}

	streamKeyIdPairs := [][]string{}

	streamsKeyIdsOffset := 1

	if isBlocking {
		streamsKeyIdsOffset = 3
	}

	for i := range pairsCount {
		storeKey, ok := data.cmd.ArgString(i + streamsKeyIdsOffset)
		if !ok {
			return &resp.Error{Msg: "ERR invalid stream name value for XREAD command"}
		}

		streamId, ok := data.cmd.ArgString(i + pairsCount + streamsKeyIdsOffset)
		if !ok {
			return &resp.Error{Msg: "ERR invalid stream id value for XREAD command"}
		}

		streamKeyIdPairs = append(streamKeyIdPairs, []string{storeKey, streamId})
	}

	streams, err := data.store.Xread(streamKeyIdPairs, blockingTimeoutMs, isBlocking)
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
