package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	internalStore "github.com/codecrafters-io/redis-starter-go/internal/store"
)

func handleSet(cmd *Command, store *store.Store) resp.Value {
	argsLen := cmd.ArgsLen()
	if argsLen < 2 || argsLen > 4 {
		return &resp.Error{Msg: "ERR wrong number of arguments for SET command"}
	}

	key, ok := cmd.ArgString(0)
	if !ok {
		return &resp.Error{Msg: "ERR invalid key value for SET command"}
	}

	value, ok := cmd.ArgBytes(1)
	if !ok {
		return &resp.Error{Msg: "ERR invalid value for SET command"}
	}

	var expiryType internalStore.ExpiryType
	var expTime int

	if argsLen > 2 && cmd.Args[2] != nil {
		expValue, ok := cmd.ArgString(2)

		if !ok {
			return &resp.Error{Msg: "ERR invalid EXP value"}
		}

		expiryType, ok = internalStore.ProcessExpType(expValue)

		if !ok {
			return &resp.Error{Msg: "ERR invalid EXP value"}
		}
	}

	if argsLen > 3 && cmd.Args[3] != nil {
		expTime, ok = cmd.ArgInt(3)

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

	if ok := store.Set(key, value, expiryType, expTime); !ok {
		return &resp.Error{Msg: "ERR during executing store SET command"}
	}

	return &resp.SimpleString{Bytes: []byte("OK")}
}
