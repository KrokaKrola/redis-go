package server

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/commands"
	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type RedisServer struct {
	port     uint16
	listener net.Listener
}

func NewRedisServer(port uint16) *RedisServer {
	return &RedisServer{
		port: port,
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
			logger.Error("error accepting connection", "err", err)
		}

		go r.handleConnection(conn)
	}
}

func (r *RedisServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	logger.Info("accepted new connection")
	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)
	decoder := resp.NewDecoder(reader)
	encoder := resp.NewEncoder(writer)
	defer writer.Flush()

	for {
		value, derr := decoder.Read()

		if derr != nil {
			if errors.Is(derr, io.EOF) {
				break
			}

			encoder.Write(&resp.Error{Msg: "ERR protocol error"})
			writer.Flush()
			continue
		}

		cmd, perr := commands.Parse(value)

		if perr != nil {
			encoder.Write(perr)
		} else {
			out := commands.Dispatch(cmd)
			encoder.Write(out)
		}

		writer.Flush()
	}
}
