package server

import (
	"bufio"
	_ "embed"
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

//go:embed empty.rdb
var emptyRDB []byte

type Session struct {
	conn                 net.Conn
	transactions         *transactions.Transactions
	id                   string
	decoder              *resp.Decoder
	encoder              *resp.Encoder
	writer               *bufio.Writer
	reader               *bufio.Reader
	serverCtx            *commands.ServerContext
	isReplicationSession bool
	countingReader       *resp.CountingReader
}

var nextClientId int64

func NewSession(conn net.Conn, store *store.Store, transactions *transactions.Transactions, isReplica bool, replicasRegistry *ReplicasRegistry, replicationId string, isReplicationSession bool) *Session {
	id := fmt.Sprintf("%d-%s", atomic.AddInt64(&nextClientId, 1), conn.RemoteAddr().String())

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	cr := &resp.CountingReader{
		R: reader,
	}
	decoder := resp.NewDecoder(cr)
	encoder := resp.NewEncoder(writer)

	return &Session{
		conn:                 conn,
		transactions:         transactions,
		id:                   id,
		decoder:              decoder,
		encoder:              encoder,
		writer:               writer,
		reader:               reader,
		isReplicationSession: isReplicationSession,
		serverCtx: &commands.ServerContext{
			IsReplica:        isReplica,
			ReplicasRegistry: replicasRegistry,
			Store:            store,
			ReplicationId:    replicationId,
		},
		countingReader: cr,
	}
}

func (s *Session) getRemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

func (s *Session) handleDecoderError(err error) (stop bool) {
	// Stop on EOF or any network connection error (closed, reset, etc.)
	if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) {
		s.transactions.Discard(s.id)
		return true
	}

	// Also stop on any net.OpError (covers "use of closed network connection")
	var opErr *net.OpError
	if errors.As(err, &opErr) {
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
		offset := s.countingReader.Count

		if s.isReplicationSession {
			logger.Debug("read offset", "offset", offset)
		}

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

		if s.isReplicationSession {
			s.serverCtx.ReplicationOffset = offset
		}

		out := s.executeCommand(cmd)

		logger.Debug("executeCommands result", slog.Any("out", out))

		if cmd.Name == commands.SET_COMMAND && !s.serverCtx.IsReplica {
			logger.Debug("replicas broadcasting", slog.Any("value", value))
			s.serverCtx.ReplicasRegistry.BroadcastRespValue(value)
			s.serverCtx.MasterOffset += resp.Size(value)
		}

		// no-op case, continue
		if out == nil {
			s.writer.Flush()
			continue
		}

		if s.isReplicationSession && cmd.Name != commands.REPLCONF {
			continue
		}

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

	if cmd.Name == commands.PSYNC && !s.serverCtx.IsReplica {
		return s.handlePsync(cmd)
	}

	if s.transactions.IsActive(s.id) {
		if err := s.transactions.Queue(s.id, cmd); err != nil {
			return &resp.Error{Msg: err.Error()}
		}

		return &resp.SimpleString{Bytes: []byte("QUEUED")}
	}

	handlerContext := &commands.HandlerContext{
		Cmd:        cmd,
		RemoteAddr: s.getRemoteAddr(),
	}

	return commands.Dispatch(s.serverCtx, handlerContext)
}

func (s *Session) handlePsync(_ *commands.Command) resp.Value {
	_, ok := s.serverCtx.ReplicasRegistry.GetReplica(s.getRemoteAddr())
	if !ok {
		return &resp.Error{Msg: "ERR replica handshake failed"}
	}

	replId := s.serverCtx.ReplicationId
	offset := 0

	if err := s.encoder.Write(&resp.SimpleString{
		Bytes: fmt.Appendf(nil, "%s %s %d", "FULLRESYNC", replId, offset),
	}); err != nil {
		return &resp.Error{
			Msg: "ERR sending FULLRESYNC response to replica",
		}
	}

	s.writer.Flush()

	if _, err := fmt.Fprintf(s.writer, "$%d\r\n%s", len(emptyRDB), emptyRDB); err != nil {
		logger.Error("failed to send RDB to replica", "error", err)
		s.conn.Close()
		return nil
	}

	s.writer.Flush()

	if err := s.serverCtx.ReplicasRegistry.AddReplicaConnection(s.conn); err != nil {
		logger.Error("failed to add replicas connection", "error", err)
		s.conn.Close()
		return nil
	}

	return nil
}

func (s *Session) handleMulti(cmd *commands.Command) resp.Value {
	if !s.transactions.IsActive(s.id) {
		handlerContext := &commands.HandlerContext{
			Cmd:        cmd,
			RemoteAddr: s.getRemoteAddr(),
		}

		out := commands.Dispatch(s.serverCtx, handlerContext)

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
		handlerContext := &commands.HandlerContext{
			Cmd:        c,
			RemoteAddr: s.getRemoteAddr(),
		}
		return commands.Dispatch(s.serverCtx, handlerContext)
	})
}

func (s *Session) handleDiscard(cmd *commands.Command) resp.Value {
	if !s.transactions.IsActive(s.id) {
		return &resp.Error{Msg: "ERR DISCARD without MULTI"}
	}

	s.transactions.Discard(s.id)

	handlerContext := &commands.HandlerContext{
		Cmd:        cmd,
		RemoteAddr: s.getRemoteAddr(),
	}

	return commands.Dispatch(s.serverCtx, handlerContext)
}
