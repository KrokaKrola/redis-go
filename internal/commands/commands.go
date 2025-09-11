package commands

import (
	"bytes"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Name string

const (
	PING_COMMAND Name = "PING"
	ECHO_COMMAND Name = "ECHO"
)

type Command struct {
	Name Name
	Args []resp.Value
}

func Parse(v resp.Value) (*Command, resp.Value) {
	cmd := &Command{}

	switch v := v.(type) {
	case *resp.Array:
		if cerr := cmd.processArray(v); cerr != nil {
			return nil, cerr
		}

		return cmd, nil
	default:
		return nil, &resp.Error{Msg: fmt.Sprintf("ERR unknown data type: %+v", v)}
	}
}

func Dispatch(cmd *Command) resp.Value {
	switch cmd.Name {
	case PING_COMMAND:
		if len(cmd.Args) == 0 {
			return &resp.SimpleString{S: []byte("PONG")}
		}

		if len(cmd.Args) == 1 {
			b, ok := valueAsBytes(cmd.Args[0])
			if !ok {
				return &resp.Error{Msg: "ERR invalid argument for ECHO command"}
			}

			return &resp.BulkString{B: b}
		}

		return &resp.Error{Msg: "ERR Invalid arguments for PING command"}
	case ECHO_COMMAND:
		if len(cmd.Args) != 1 {
			return &resp.Error{Msg: "ERR wrong number of arguments for ECHO command"}
		}
		b, ok := valueAsBytes(cmd.Args[0])
		if !ok {
			return &resp.Error{Msg: "ERR invalid argument for ECHO command"}
		}
		return &resp.BulkString{B: b}
	default:
		return &resp.Error{Msg: fmt.Sprintf("ERR unknown command name: %s", cmd.Name)}
	}
}

func (c *Command) processArray(arr *resp.Array) resp.Value {
	if arr.Null || len(arr.Elems) == 0 {
		return &resp.Error{Msg: "ERR invalid size of array"}
	}

	b, ok := valueAsBytes(arr.Elems[0])
	if !ok {
		return &resp.Error{Msg: "ERR protocol error"}
	}
	name := getCommandName(b)

	if name == "" {
		return &resp.Error{Msg: "ERR unknown command"}
	}
	c.Name = name

	if len(arr.Elems) > 1 {
		c.Args = append(c.Args, arr.Elems[1:]...)
	}

	return nil
}

func getCommandName(name []byte) Name {
	if bytes.EqualFold(name, []byte(PING_COMMAND)) {
		return PING_COMMAND
	} else if bytes.EqualFold(name, []byte(ECHO_COMMAND)) {
		return ECHO_COMMAND
	} else {
		return ""
	}
}

func valueAsBytes(v resp.Value) ([]byte, bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return nil, false
		}
		return x.B, true
	case *resp.SimpleString:
		return x.S, true
	default:
		return nil, false
	}
}
