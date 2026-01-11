package server

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
	"github.com/codecrafters-io/redis-starter-go/internal/transactions"
)

const shutdownTimeout = 5 * time.Second

type RedisServer struct {
	port             int
	listener         net.Listener
	store            *store.Store
	transactions     *transactions.Transactions
	wg               sync.WaitGroup // tracks active connections
	isReplica        bool
	replicasRegistry *ReplicasRegistry
}

func NewRedisServer(port int, isReplica bool) *RedisServer {
	return &RedisServer{
		port:             port,
		store:            store.NewStore(),
		transactions:     transactions.NewTransactions(),
		isReplica:        isReplica,
		replicasRegistry: NewReplicasRegistry(),
	}
}

func (r *RedisServer) ConnectToMaster(replicaOf string, replicaPort int) error {
	parts := strings.Split(replicaOf, " ")
	conn, err := net.Dial("tcp", net.JoinHostPort(parts[0], parts[1]))

	if err != nil {
		return err
	}

	encoder := resp.NewEncoder(conn)
	decoder := resp.NewDecoder(bufio.NewReader(conn))

	pingMsg := &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{
				Bytes: []byte("PING"),
			},
		},
	}

	if err = encoder.Write(pingMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending PING request to master %s from replica", replicaOf), err)
	}

	if _, err := decoder.Read(); err != nil {
		return errors.Join(fmt.Errorf("error reading PING response from replica %s", replicaOf), err)
	}

	replConfMsg := &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{Bytes: []byte("REPLCONF")},
			&resp.BulkString{Bytes: []byte("listening-port")},
			&resp.BulkString{Bytes: fmt.Append(nil, replicaPort)},
		},
	}

	if err = encoder.Write(replConfMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending REPLCONF listening-port request to master %s from replica", replicaOf), err)
	}

	if _, err := decoder.Read(); err != nil {
		return errors.Join(fmt.Errorf("error reading REPLCONF listenin-port response from replica %s", replicaOf), err)
	}

	replConfMsg = &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{Bytes: []byte("REPLCONF")},
			&resp.BulkString{Bytes: []byte("capa")},
			&resp.BulkString{Bytes: []byte("psync2")},
		},
	}

	if err = encoder.Write(replConfMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending REPLCONF capa request to master %s from replica", replicaOf), err)
	}

	if _, err := decoder.Read(); err != nil {
		return errors.Join(fmt.Errorf("error reading REPLCONF capa response from replica %s", replicaOf), err)
	}

	return nil
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

	session := NewSession(conn, r.store, r.transactions, r.isReplica, r.replicasRegistry)
	session.Run()
}
