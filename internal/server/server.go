package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

type RedisServer struct {
	port     uint16
	listener net.Listener
	store    *store.Store
}

func NewRedisServer(port uint16) *RedisServer {
	return &RedisServer{
		port:  port,
		store: store.NewStore(),
	}
}

func (r *RedisServer) Listen() error {
	if r.port == 0 {
		return fmt.Errorf("port is not specified")
	}

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", r.port))
	if err != nil {
		return err
	}

	defer l.Close()

	logger.Info("Started server",
		"address", l.Addr(),
		"port", r.port,
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	r.listener = l
	go r.acceptConnections()
	<-sigChan

	return nil
}

func (r *RedisServer) acceptConnections() {
	for {
		conn, err := r.listener.Accept()

		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}

			logger.Error("error accepting connection", "err", err)
			continue
		}

		go r.handleConnection(conn)
	}
}

func (r *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	logger.Debug("accepted new connection", "RemoteAddr", conn.RemoteAddr())
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	decoder := resp.NewDecoder(reader)
	encoder := resp.NewEncoder(writer)
	defer writer.Flush()

	for {
		value, derr := decoder.Read()

		logger.Debug("decoder.Read result", slog.Any("value", value), slog.Any("derr", derr))

		if derr != nil {
			if errors.Is(derr, io.EOF) {
				break
			}

			encoder.Write(&resp.Error{Msg: "ERR protocol error"})
			writer.Flush()
			continue
		}

		cmd, perr := commands.Parse(value)

		logger.Debug("commands.Parse result", slog.Any("cmd", cmd), slog.Any("perr", perr))

		if perr != nil {
			encoder.Write(&resp.Error{Msg: perr.Error()})
		} else {
			out := commands.Dispatch(cmd, r.store)
			logger.Debug("commands.Dispatch result", slog.Any("out", out))

			if err := encoder.Write(out); err != nil {
				encoder.Write(&resp.Error{Msg: fmt.Sprintf("ERR encoder failed to write a response: %T", err.Error())})
			}

		}

		writer.Flush()
	}
}
