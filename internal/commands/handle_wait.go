package commands

import (
	"bufio"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func handleWait(serverCtx *ServerContext, handlerCtx *HandlerContext) resp.Value {
	argsLen := handlerCtx.Cmd.ArgsLen()

	if argsLen < 2 {
		return &resp.Error{Msg: "ERR wrong number of arguments for 'wait' command"}
	}

	numReplicas, ok := handlerCtx.Cmd.ArgInt(0)
	if !ok {
		return &resp.Error{Msg: "ERR wrong number of arguments for 'wait' command"}
	}

	timeout, ok := handlerCtx.Cmd.ArgInt(1)
	if !ok {
		return &resp.Error{Msg: "ERR wrong number of arguments for 'wait' command"}
	}

	timeoutDuration := time.Duration(timeout) * time.Millisecond

	replicas := serverCtx.ReplicasRegistry.GetAllReplicas()

	if len(replicas) == 0 {
		return &resp.Integer{Number: 0}
	}

	if serverCtx.MasterOffset == 0 {
		return &resp.Integer{Number: int64(len(replicas))}
	}

	masterOffset := serverCtx.MasterOffset

	// Send REPLCONF GETACK * to all replicas
	replConfMsg := &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{Bytes: []byte("REPLCONF")},
			&resp.BulkString{Bytes: []byte("GETACK")},
			&resp.BulkString{Bytes: []byte("*")},
		},
	}

	for _, repl := range replicas {
		writer := bufio.NewWriter(repl.Connection)
		encoder := resp.NewEncoder(writer)
		if err := encoder.Write(replConfMsg); err != nil {
			logger.Error("ERR writing GETACK to replica", "error", err)
			continue
		}
		writer.Flush()
	}

	// Poll for acknowledgments until timeout or enough replicas respond
	deadline := time.Now().Add(timeoutDuration)
	pollInterval := 10 * time.Millisecond

	for time.Now().Before(deadline) {
		count := 0
		for _, repl := range replicas {
			if repl.AckOffset >= masterOffset {
				count++
			}
		}

		if count >= numReplicas {
			return &resp.Integer{Number: int64(count)}
		}

		time.Sleep(pollInterval)
	}

	// Final count after timeout
	count := 0
	for _, repl := range replicas {
		if repl.AckOffset >= masterOffset {
			count++
		}
	}

	return &resp.Integer{Number: int64(count)}
}
