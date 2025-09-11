package commands

import (
	"fmt"
	"strings"

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
	case *resp.BulkString:
		cmd.processBulkString(v)
		return cmd, nil
	default:
		return nil, resp.Error{Msg: fmt.Sprintf("Unknown data type: %+v", v)}
	}
}

func Dispatch(cmd *Command) resp.Value {
	switch cmd.Name {
	case PING_COMMAND:
		return &resp.SimpleString{S: []byte("PONG")}
	case ECHO_COMMAND:
		if len(cmd.Args) == 0 {
			return &resp.Error{Msg: "Invalid number of arguments for ECHO command"}
		}
		bs, ok := cmd.Args[0].(*resp.BulkString)
		if !ok {
			return &resp.Error{Msg: "Invalid arguments for ECHO command"}
		}

		return bs
	default:
		return &resp.Error{Msg: fmt.Sprintf("Unknown command name: %s", cmd.Name)}
	}
}

func (c *Command) processArray(arr *resp.Array) resp.Value {
	// TODO: process arr.Null
	// TODO: process empty array

	var nameElem *resp.BulkString
	it := 0
	var args []resp.Value

	for it < len(arr.Elems) {
		if it == 0 {
			bs, ok := arr.Elems[it].(*resp.BulkString)
			if !ok {
				return resp.Error{Msg: "Expected BulkString for command name"}
			}
			nameElem = bs
		} else {
			args = append(args, arr.Elems[it])
		}
		it++
	}

	switch strings.ToUpper(string(nameElem.B)) {
	case string(PING_COMMAND):
		c.Name = PING_COMMAND
	case string(ECHO_COMMAND):
		c.Name = ECHO_COMMAND
	default:
		return resp.Error{Msg: "Unknown command name"}
	}

	c.Args = args

	return nil
}

func (c *Command) processBulkString(str *resp.BulkString) resp.Value {
	return nil
}
