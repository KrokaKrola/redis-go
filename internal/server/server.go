package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/transactions"
)

type RedisServer struct {
	port         uint16
	listener     net.Listener
	store        *store.Store
	transactions *transactions.Transactions
}

func NewRedisServer(port uint16) *RedisServer {
	return &RedisServer{
		port:         port,
		store:        store.NewStore(),
		transactions: transactions.NewTransactions(),
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
	session := NewSession(conn, r.store, r.transactions)
	session.Run()
}
