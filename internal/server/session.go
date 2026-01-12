package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"sync/atomic"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/transactions"
)

type Session struct {
	conn             net.Conn
	store            *store.Store
	transactions     *transactions.Transactions
	id               string
	decoder          *resp.Decoder
	encoder          *resp.Encoder
	writer           *bufio.Writer
	reader           *bufio.Reader
	isReplica        bool
	replicasRegistry *ReplicasRegistry
	replicationId    string
}

var nextClientId int64

func NewSession(conn net.Conn, store *store.Store, transactions *transactions.Transactions, isReplica bool, replicasRegistry *ReplicasRegistry, replicationId string) *Session {
	id := fmt.Sprintf("%d-%s", atomic.AddInt64(&nextClientId, 1), conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	decoder := resp.NewDecoder(reader)
	encoder := resp.NewEncoder(writer)

	return &Session{
		conn:             conn,
		store:            store,
		transactions:     transactions,
		id:               id,
		decoder:          decoder,
		encoder:          encoder,
		writer:           writer,
		reader:           reader,
		isReplica:        isReplica,
		replicasRegistry: replicasRegistry,
		replicationId:    replicationId,
	}
}

func (s *Session) getRemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Session) handleDecoderError(err error) (stop bool) {
	if errors.Is(err, io.EOF) {
		s.transactions.Discard(s.id)
		return true
	}

	s.encoder.Write(&resp.Error{Msg: "ERR protocol error"})
	s.writer.Flush()

	return false
}

func (s *Session) writeError(msg string) {
	s.encoder.Write(&resp.Error{Msg: msg})
	s.writer.Flush()
}

func (s *Session) Run() {
	logger.Debug("accepted new connection", "RemoteAddr", s.conn.RemoteAddr())

	defer s.conn.Close()
	defer s.writer.Flush()

	for {
		value, derr := s.decoder.Read()

		logger.Debug("decoder.Read result", slog.Any("value", value), slog.Any("derr", derr))

		if derr != nil {
			if stop := s.handleDecoderError(derr); stop {
				break
			}

			continue
		}

		cmd, perr := commands.Parse(value)

		logger.Debug("commands.Parse result", slog.Any("cmd", cmd), slog.Any("perr", perr))

		if perr != nil {
			s.writeError(perr.Error())
			continue
		}

		out := s.executeCommand(cmd)

		err := s.encoder.Write(out)
		if err != nil {
			s.encoder.Write(&resp.Error{Msg: fmt.Sprintf("ERR encoder failed to write a response: %s", err.Error())})
		}

		s.writer.Flush()
	}
}

func (s *Session) executeCommand(cmd *commands.Command) resp.Value {
	if cmd.Name == commands.MULTI_COMMAND {
		return s.handleMulti(cmd)
	}

	if cmd.Name == commands.EXEC_COMMAND {
		return s.handleExec()
	}

	if cmd.Name == commands.DISCARD_COMMAND {
		return s.handleDiscard(cmd)
	}

	if s.transactions.IsActive(s.id) {
		if err := s.transactions.Queue(s.id, cmd); err != nil {
			return &resp.Error{Msg: err.Error()}
		}

		return &resp.SimpleString{Bytes: []byte("QUEUED")}
	}

	serverContext := &commands.ServerContext{
		IsReplica:        s.isReplica,
		ReplicasRegistry: s.replicasRegistry,
		Store:            s.store,
		ReplicationId:    s.replicationId,
	}
	handlerContext := &commands.HandlerContext{
		Cmd:        cmd,
		RemoteAddr: s.getRemoteAddr(),
	}

	return commands.Dispatch(serverContext, handlerContext)
}

func (s *Session) handleMulti(cmd *commands.Command) resp.Value {
	if !s.transactions.IsActive(s.id) {
		serverContext := &commands.ServerContext{
			IsReplica:        s.isReplica,
			ReplicasRegistry: s.replicasRegistry,
			Store:            s.store,
			ReplicationId:    s.replicationId,
		}
		handlerContext := &commands.HandlerContext{
			Cmd:        cmd,
			RemoteAddr: s.getRemoteAddr(),
		}

		out := commands.Dispatch(serverContext, handlerContext)

		switch out.(type) {
		case *resp.Error:
			return out
		default:
			s.transactions.Begin(s.id)
			return out
		}
	} else {
		return &resp.Error{Msg: "ERR multi calls cannot be nested"}
	}
}

func (s *Session) handleExec() resp.Value {
	if !s.transactions.IsActive(s.id) {
		return &resp.Error{Msg: "ERR EXEC without MULTI"}
	}

	return s.transactions.ExecuteAndDiscard(s.id, func(c *commands.Command) resp.Value {
		serverContext := &commands.ServerContext{
			IsReplica:        s.isReplica,
			ReplicasRegistry: s.replicasRegistry,
			Store:            s.store,
			ReplicationId:    s.replicationId,
		}
		handlerContext := &commands.HandlerContext{
			Cmd:        c,
			RemoteAddr: s.getRemoteAddr(),
		}
		return commands.Dispatch(serverContext, handlerContext)
	})
}

func (s *Session) handleDiscard(cmd *commands.Command) resp.Value {
	if !s.transactions.IsActive(s.id) {
		return &resp.Error{Msg: "ERR DISCARD without MULTI"}
	}

	s.transactions.Discard(s.id)

	serverContext := &commands.ServerContext{
		IsReplica:        s.isReplica,
		ReplicasRegistry: s.replicasRegistry,
		Store:            s.store,
		ReplicationId:    s.replicationId,
	}
	handlerContext := &commands.HandlerContext{
		Cmd:        cmd,
		RemoteAddr: s.getRemoteAddr(),
	}

	return commands.Dispatch(serverContext, handlerContext)
}
