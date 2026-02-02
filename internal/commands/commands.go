package commands

import (
	"fmt"

	"github.com/codecrafters-io/redis-starter-go/internal/replica"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type Command struct {
	Name Name
	Args []resp.Value
}

func (c *Command) ArgsLen() int {
	return len(c.Args)
}

func (c *Command) ArgBytes(idx int) ([]byte, bool) {
	if idx >= c.ArgsLen() {
		return []byte{}, false
	}
	return valueAsBytes(c.Args[idx])
}

func (c *Command) ArgString(idx int) (string, bool) {
	if idx >= c.ArgsLen() {
		return "", false
	}
	return valueAsString(c.Args[idx])
}

func (c *Command) ArgInt(idx int) (int, bool) {
	if idx >= c.ArgsLen() {
		return 0, false
	}
	return valueAsInteger(c.Args[idx])
}

func (c *Command) ArgFloat(idx int) (float64, bool) {
	if idx >= c.ArgsLen() {
		return 0, false
	}
	return valueAsFloat(c.Args[idx])
}

func Parse(v resp.Value) (*Command, error) {
	switch v := v.(type) {
	case *resp.Array:
		cmd, cerr := newCommandFromRespArray(v)
		if cerr != nil {
			return nil, cerr
		}

		return cmd, nil
	default:
		return nil, fmt.Errorf("ERR got unknown data type during parsing: %+v", v)
	}
}

func newCommandFromRespArray(arr *resp.Array) (*Command, error) {
	cmd := &Command{}

	if arr.Null || len(arr.Elements) == 0 {
		return nil, fmt.Errorf("ERR invalid size of array")
	}

	b, ok := valueAsBytes(arr.Elements[0])
	if !ok {
		return nil, fmt.Errorf("ERR protocol error")
	}
	name := getCommandName(b)

	if name == "" {
		return nil, fmt.Errorf("ERR unknown command name -> %s", b)
	}

	cmd.Name = name

	if len(arr.Elements) > 1 {
		cmd.Args = append(cmd.Args, arr.Elements[1:]...)
	}

	return cmd, nil

}

type ServerContext struct {
	IsReplica         bool
	ReplicasRegistry  replica.ReplicasRegistry
	Store             *store.Store
	ReplicationId     string
	ReplicationOffset int
}

type HandlerContext struct {
	Cmd        *Command
	RemoteAddr string
}

type handlerFn func(*ServerContext, *HandlerContext) resp.Value

var handlers = map[Name]handlerFn{
	PING_COMMAND:    handlePing,
	ECHO_COMMAND:    handleEcho,
	GET_COMMAND:     handleGet,
	SET_COMMAND:     handleSet,
	RPUSH_COMMAND:   handlePush,
	LPUSH_COMMAND:   handlePush,
	LRANGE_COMMAND:  handleLrange,
	LLEN_COMMAND:    handleLlen,
	LPOP_COMMAND:    handleLpop,
	BLPOP_COMMAND:   handleBlpop,
	TYPE_COMMAND:    handleType,
	XADD_COMMAND:    handleXadd,
	XRANGE_COMMAND:  handleXrange,
	XREAD_COMMAND:   handleXread,
	INCR_COMMAND:    handleIncr,
	MULTI_COMMAND:   handleMulti,
	DISCARD_COMMAND: handleDiscard,
	INFO_COMMAND:    handleInfo,
	REPLCONF:        handleReplconf,
}

func Dispatch(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	if handler, ok := handlers[handlerCtx.Cmd.Name]; ok {
		return handler(serverCtx, handlerCtx)
	}

	return &resp.Error{Msg: fmt.Sprintf("ERR handler for %s is not implemented", handlerCtx.Cmd.Name)}
}
