package server

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
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
	log.Println("accepted new connection")

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)

		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Println("EOF of the conn")

				break
			}

			log.Fatal("error while reading connection happened: ", err)
		}

		log.Printf("read %d bytes from conn\n", n)

		if n == 0 {
			continue
		}

		conn.Write([]byte("+PONG\r\n"))
	}
}
