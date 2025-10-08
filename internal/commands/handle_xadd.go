package commands

import (
	"strconv"
	"strings"

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

	sId, ok := cmd.ArgString(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid stream-id value for XADD command"}
	}

	streamId, ok := parseStreamId(sId)

	if !ok {
		return &resp.Error{Msg: "ERR invalid stream-id value for XADD command"}
	}

	// if manual ID check the seq number and ms time parts
	if !streamId.AutoSeq && !streamId.AutoFull {
		if streamId.MsTime == 0 && streamId.Seq == 0 {
			return &resp.Error{Msg: "ERR The ID specified in XADD must be greater than 0-0"}
		}
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

	id, err := store.Xadd(key, streamId, fields)

	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	return &resp.BulkString{Bytes: []byte(id)}
}

func parseStreamId(id string) (streamId store.StreamIdSpec, ok bool) {
	if id == "*" {
		return store.StreamIdSpec{
			AutoFull: true,
		}, true
	}

	before, after, found := strings.Cut(id, "-")

	if !found {
		return store.StreamIdSpec{}, false
	}

	msTime, err := strconv.ParseUint(before, 10, 64)
	if err != nil {
		return store.StreamIdSpec{}, false
	}

	if after == "*" {
		return store.StreamIdSpec{
			MsTime:  msTime,
			AutoSeq: true,
		}, true
	}

	seq, err := strconv.ParseUint(after, 10, 64)
	if err != nil {
		return store.StreamIdSpec{}, false
	}

	return store.StreamIdSpec{MsTime: msTime, Seq: seq}, true
}
