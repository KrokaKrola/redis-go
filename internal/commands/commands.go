package commands

import (
	"bytes"
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
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
		argsLen := len(cmd.Args)
		if argsLen < 2 || argsLen > 4 {
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

		var expiryType store.ExpiryType
		var expTime int

		if argsLen > 2 && cmd.Args[2] != nil {
			expValue, ok := valueAsString(cmd.Args[2])

			if !ok {
				return &resp.Error{Msg: "ERR invalid EXP value"}
			}

			expiryType, ok = store.ProcessExpType(expValue)

			if !ok {
				return &resp.Error{Msg: "ERR invalid EXP value"}
			}
		}

		if argsLen > 3 && cmd.Args[3] != nil {
			expTime, ok = valueAsInteger(cmd.Args[3])

			if !ok {
				return &resp.Error{Msg: "ERR invalid expTime for SET command"}
			}

			if expTime <= 0 {
				return &resp.Error{Msg: "ERR invalid expTime for SET command"}
			}
		}

		if expiryType != "" && expTime == 0 {
			return &resp.Error{Msg: "ERR invalid expTime for SET command"}
		}

		if ok := s.Set(key, value, expiryType, expTime); !ok {
			return &resp.Error{Msg: "ERR during executing store SET command"}
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
	case RPUSH_COMMAND:
		argsLen := len(cmd.Args)
		if argsLen < 2 {
			return &resp.Error{Msg: "ERR wrong number of arguments for RPUSH command"}
		}

		key, ok := valueAsString(cmd.Args[0])
		if !ok {
			return &resp.Error{Msg: "ERR invalid key value for RPUSH command"}
		}

		var values []string

		for argsLen-len(values) != 1 {
			value, ok := valueAsString(cmd.Args[len(values)+1])
			if !ok {
				return &resp.Error{Msg: "ERR invalid type of RPUSH arguments list item"}
			}

			values = append(values, value)
		}

		if len(values) == 0 {
			return &resp.Error{Msg: "ERR empty values for RPUSH command"}
		}

		v, ok := s.Rpush(key, store.List{L: values})

		if !ok {
			return &resp.Error{Msg: "WRONGTYPE Operation against a key holding the wrong kind of value"}
		}

		return &resp.Integer{N: v}
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
		return fmt.Errorf("ERR something went wrong while getting command name, probably command name is not defined")
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
	} else if bytes.EqualFold(name, []byte(RPUSH_COMMAND)) {
		return RPUSH_COMMAND
	} else {
		return ""
	}
}
