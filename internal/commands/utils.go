package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

func valueAsBytes(v resp.Value) (value []byte, ok bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return nil, false
		}

		return x.Bytes, true
	case *resp.SimpleString:
		return x.Bytes, true
	default:
		return nil, false
	}
}

func valueAsString(v resp.Value) (value string, ok bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return "", false
		}

		return string(x.Bytes), true
	case *resp.SimpleString:
		return string(x.Bytes), true
	default:
		return "", false
	}
}

func valueAsInteger(v resp.Value) (value int, ok bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return 0, false
		}

		v, err := strconv.Atoi(string(x.Bytes))

		if err != nil {
			return 0, false
		}

		return v, true
	case *resp.SimpleString:
		v, err := strconv.Atoi(string(x.Bytes))

		if err != nil {
			return 0, false
		}

		return v, true
	default:
		return 0, false
	}
}

func valueAsFloat(v resp.Value) (value float64, ok bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return 0, false
		}

		v, err := strconv.ParseFloat(string(x.Bytes), 64)

		if err != nil {
			return 0, false
		}

		return v, true
	case *resp.SimpleString:
		v, err := strconv.ParseFloat(string(x.Bytes), 64)

		if err != nil {
			return 0, false
		}

		return v, true
	default:
		return 0, false
	}
}

func populateRespArrayFromStream(stream store.Stream) *resp.Array {
	arr := &resp.Array{}

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
