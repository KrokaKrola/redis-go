package commands

import (
	"bufio"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/replica"
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
	resChan := make(chan resp.Value)

	replicas := serverCtx.ReplicasRegistry.GetAllReplicas()

	if len(replicas) == 0 {
		return &resp.Integer{Number: 0}
	}

	if serverCtx.ReplicationOffset == 0 {
		return &resp.Integer{Number: int64(len(replicas))}
	}

	for _, repl := range replicas {
		go func(repl *replica.Replica) {
			writer := bufio.NewWriter(repl.Connection)
			reader := bufio.NewReader(repl.Connection)
			encoder := resp.NewEncoder(writer)
			decoder := resp.NewDecoder(reader)
			replConfMsg := &resp.Array{
				Elements: []resp.Value{
					&resp.BulkString{Bytes: []byte("REPLCONF")},
					&resp.BulkString{Bytes: []byte("GETACK")},
					&resp.BulkString{Bytes: []byte("*")},
				},
			}
			if err := encoder.Write(replConfMsg); err != nil {
				logger.Error("ERR wrong writing replica configuration", "error", err)
				return
			}

			writer.Flush()

			response, err := decoder.Read()
			if err != nil {
				return
			}

			resChan <- response
		}(repl)
	}

	count := 0
	timer := time.After(timeoutDuration)

	for {
		select {
		case <-resChan:
			count++
			if count >= numReplicas {
				return &resp.Integer{Number: int64(count)}
			}

		case <-timer:
			return &resp.Integer{Number: int64(count)}
		}
	}
}
