package commands

import (
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
	case resp.Array:
		cmd.processArray(v)
		return cmd, nil
	case resp.BulkString:
		cmd.processBulkString(v)
		return cmd, nil
	default:
		return nil, resp.Error{Msg: fmt.Sprintf("Unknown data type: %v", v)}
	}
}

func Dispatch(cmd *Command) resp.Value {
	return nil
}

func (c *Command) processArray(arr resp.Array) error {
	// TODO: process arr.Null
	return nil
}

func (c *Command) processBulkString(str resp.BulkString) error {
	return nil
}
