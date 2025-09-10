package commands

import "github.com/codecrafters-io/redis-starter-go/internal/resp"

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
	return nil, nil
}

func Dispatch(cmd *Command) resp.Value {
	return nil
}
