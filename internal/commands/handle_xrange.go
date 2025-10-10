package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleXrange(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()

	if argsLen < 3 {
		return &resp.Error{Msg: "ERR invalid number of arguments for XRANGE command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for XRANGE command"}
	}

	start, ok := cmd.ArgString(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid start value for XRANGE command"}
	}

	end, ok := cmd.ArgString(2)
	if !ok {
		return &resp.Error{Msg: "ERR invalid end value for XRANGE command"}
	}

	stream, err := store.Xrange(key, start, end)
	if err != nil {
		return &resp.Error{Msg: err.Error()}
	}

	arr := &resp.Array{}

	if len(stream.Elements) == 0 {
		return arr
	}

	for _, el := range stream.Elements {
		fields := &resp.Array{}

		for _, field := range el.Fields {
			for i := range 2 {
				fields.Elements = append(fields.Elements, &resp.BulkString{Bytes: []byte(field[i])})
			}
		}

		arr.Elements = append(arr.Elements, resp.Value(&resp.Array{
			Elements: []resp.Value{
				&resp.BulkString{Bytes: fmt.Appendf(nil, "%d-%d", el.Id.MsTime, el.Id.Seq)},
				fields,
			},
		},
		))
	}

	return arr
}
