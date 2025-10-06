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

	streamEntryId, ok := cmd.ArgString(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid stream-id value for XADD command"}
	}

	msTime, seqNumber, isAutogenSeqNumber, isAutogen, ok := parseStreamId(streamEntryId)

	if !ok {
		return &resp.Error{Msg: "ERR invalid stream-id value for XADD command"}
	}

	// if manual ID check the seq number and ms time parts
	if !isAutogenSeqNumber && !isAutogen {
		if msTime == 0 && seqNumber == 0 {
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

	id, err := store.Xadd(key, msTime, seqNumber, isAutogenSeqNumber, isAutogen, fields)

	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	return &resp.BulkString{Bytes: []byte(id)}
}

func parseStreamId(id string) (msTime uint64, seqNumber uint64, isAutogenSeqNumber bool, isAutgen bool, ok bool) {
	if id == "*" {
		return 0, 0, false, true, true
	}

	before, after, found := strings.Cut(id, "-")

	if !found {
		return 0, 0, false, false, false
	}

	msTime, err := strconv.ParseUint(before, 10, 64)
	if err != nil {
		return 0, 0, false, false, false
	}

	if after == "*" {
		return msTime, 0, true, false, true
	}

	seqNumber, err = strconv.ParseUint(after, 10, 64)
	if err != nil {
		return 0, 0, false, false, false
	}

	return msTime, seqNumber, false, false, true
}
