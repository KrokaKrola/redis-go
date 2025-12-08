package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/transactions"
)

const shutdownTimeout = 5 * time.Second

type RedisServer struct {
	port         int
	listener     net.Listener
	store        *store.Store
	transactions *transactions.Transactions
	wg           sync.WaitGroup // tracks active connections
	isReplica    bool
}

func NewRedisServer(port int, isReplica bool) *RedisServer {
	return &RedisServer{
		port:         port,
		store:        store.NewStore(),
		transactions: transactions.NewTransactions(),
		isReplica:    isReplica,
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

	r.listener = l

	logger.Info("Started server",
		"address", l.Addr(),
		"port", r.port,
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go r.acceptConnections()

	<-sigChan
	logger.Info("Shutdown signal received, stopping server...")

	// Stop accepting new connections
	l.Close()

	// Wait for active connections to finish (with timeout)
	done := make(chan struct{})
	go func() {
		r.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Info("All connections closed gracefully")
	case <-time.After(shutdownTimeout):
		logger.Warn("Shutdown timeout reached, forcing exit")
	}

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
	r.wg.Add(1)
	defer r.wg.Done()

	session := NewSession(conn, r.store, r.transactions, r.isReplica)
	session.Run()
}
