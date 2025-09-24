package commands

import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type Name string

const (
	PING_COMMAND Name = "PING"
	ECHO_COMMAND Name = "ECHO"
	GET_COMMAND  Name = "GET"
	SET_COMMAND  Name = "SET"
)

type Command struct {
	Name Name
	Args []resp.Value
}

func Parse(v resp.Value) (*Command, error) {
	cmd := &Command{}

	switch v := v.(type) {
	case *resp.Array:
		if cerr := cmd.processArray(v); cerr != nil {
			return nil, cerr
		}

		return cmd, nil
	default:
		return nil, fmt.Errorf("ERR unknown data type: %+v", v)
	}
}

func Dispatch(cmd *Command, s *store.Store) resp.Value {
	switch cmd.Name {
	case PING_COMMAND:
		if len(cmd.Args) == 0 {
			return &resp.SimpleString{S: []byte("PONG")}
		}

		if len(cmd.Args) == 1 {
			b, ok := valueAsBytes(cmd.Args[0])
			if !ok {
				return &resp.Error{Msg: "ERR invalid argument for PING command"}
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
	case SET_COMMAND:
		if len(cmd.Args) < 2 {
			return &resp.Error{Msg: "ERR wrong number of arguments for SET command"}
		}

		key, ok := valueAsString(cmd.Args[0])
		if !ok {
			return &resp.Error{Msg: "ERR invalid key value for SET command"}
		}

		value, ok := valueAsBytes(cmd.Args[1])
		if !ok {
			return &resp.Error{Msg: "ERR invalid value for SET command"}
		}

		logger.Debug("cmd.Args[2] value", cmd.Args[2])

		var expValue string
		var expTime int

		if cmd.Args[2] != nil {
			expValue, ok = valueAsString(cmd.Args[2])

			if !ok {
				return &resp.Error{Msg: "ERR invalid value for SET command"}
			}

		}

		if cmd.Args[3] != nil {
			expTime, ok = valueAsInteger(cmd.Args[3])

			if !ok {
				return &resp.Error{Msg: "ERR invalid value for SET command"}
			}

			if expTime < 0 {
				return &resp.Error{Msg: "ERR invalid value for SET command"}
			}
		}

		if ok := s.Set(key, value, expValue, expTime); !ok {
			return &resp.Error{Msg: "ERR invalid value for SET command"}
		}

		return &resp.SimpleString{S: []byte("OK")}
	case GET_COMMAND:
		if len(cmd.Args) != 1 {
			return &resp.Error{Msg: "ERR wrong number of arguments for GET command"}
		}

		key, ok := valueAsString(cmd.Args[0])
		if !ok {
			return &resp.Error{Msg: "ERR invalid key value for GET command"}
		}

		v, ok := s.Get(key)

		if !ok {
			return &resp.BulkString{Null: true}
		}

		return &resp.BulkString{B: v}
	default:
		return &resp.Error{Msg: fmt.Sprintf("ERR unknown command name: %s", cmd.Name)}
	}
}

func (c *Command) processArray(arr *resp.Array) error {
	if arr.Null || len(arr.Elems) == 0 {
		return fmt.Errorf("ERR invalid size of array")
	}

	b, ok := valueAsBytes(arr.Elems[0])
	if !ok {
		return fmt.Errorf("ERR protocol error")
	}
	name := getCommandName(b)

	if name == "" {
		return fmt.Errorf("ERR unknown command")
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
	} else if bytes.EqualFold(name, []byte(GET_COMMAND)) {
		return GET_COMMAND
	} else if bytes.EqualFold(name, []byte(SET_COMMAND)) {
		return SET_COMMAND
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
		return nil, true
	}
}

func valueAsString(v resp.Value) (string, bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return "", false
		}

		return string(x.B), true
	case *resp.SimpleString:
		return string(x.S), true
	default:
		return "", true
	}
}

func valueAsInteger(v resp.Value) (int, bool) {
	switch x := v.(type) {
	case *resp.BulkString:
		if x.Null {
			return 0, false
		}
		v, err := strconv.Atoi(string(x.B))

		if err != nil {
			return 0, true
		}

		return v, false
	case *resp.SimpleString:
		v, err := strconv.Atoi(string(x.S))

		if err != nil {
			return 0, true
		}

		return v, false
	default:
		return 0, true
	}
}
