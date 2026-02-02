package server

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"slices"
	"strconv"
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
	replicationId    string
}

func NewRedisServer(port int, isReplica bool) *RedisServer {
	replicationId := ""

	if !isReplica {
		replicationId = "8371b4fb1155b71f4a04d3e1bc3e18c4a990aeeb"
	}

	return &RedisServer{
		port:             port,
		store:            store.NewStore(),
		transactions:     transactions.NewTransactions(),
		isReplica:        isReplica,
		replicasRegistry: NewReplicasRegistry(),
		replicationId:    replicationId,
	}
}

func (r *RedisServer) ConnectToMaster(replicaOf string, replicaPort int) error {
	parts := strings.Split(replicaOf, " ")
	conn, err := net.Dial("tcp", net.JoinHostPort(parts[0], parts[1]))

	if err != nil {
		return errors.Join(errors.New("error while trying to connect to the master server"), err)
	}

	session := NewSession(conn, r.store, r.transactions, r.isReplica, r.replicasRegistry, r.replicationId, true)

	pingMsg := &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{
				Bytes: []byte("PING"),
			},
		},
	}

	if err = session.encoder.Write(pingMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending PING request to master %s from replica", replicaOf), err)
	}

	session.writer.Flush()

	pingResponse, err := session.decoder.Read()
	if err != nil {
		return errors.Join(fmt.Errorf("error reading PING response from replica %s", replicaOf), err)
	}

	if _, ok := pingResponse.(*resp.Error); ok {
		return fmt.Errorf("invalid PING response from replica %s", replicaOf)
	}

	replConfMsg := &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{Bytes: []byte("REPLCONF")},
			&resp.BulkString{Bytes: []byte("listening-port")},
			&resp.BulkString{Bytes: fmt.Append(nil, replicaPort)},
		},
	}

	if err = session.encoder.Write(replConfMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending REPLCONF listening-port request to master %s from replica", replicaOf), err)
	}

	session.writer.Flush()

	replConfListeningPortResponse, err := session.decoder.Read()
	if err != nil {
		return errors.Join(fmt.Errorf("error reading REPLCONF listening-port response from replica %s", replicaOf), err)
	}

	if _, ok := replConfListeningPortResponse.(*resp.Error); ok {
		return fmt.Errorf("invalid REPLCONF listening-port response from replica %s", replicaOf)
	}

	replConfMsg = &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{Bytes: []byte("REPLCONF")},
			&resp.BulkString{Bytes: []byte("capa")},
			&resp.BulkString{Bytes: []byte("psync2")},
		},
	}

	if err = session.encoder.Write(replConfMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending REPLCONF capa request to master %s from replica", replicaOf), err)
	}

	session.writer.Flush()

	replConfCapaResponse, err := session.decoder.Read()
	if err != nil {
		return errors.Join(fmt.Errorf("error reading REPLCONF capa response from replica %s", replicaOf), err)
	}

	if _, ok := replConfCapaResponse.(*resp.Error); ok {
		return fmt.Errorf("invalid REPLCONF capa response from replica %s", replicaOf)
	}

	psyncInitialMsg := &resp.Array{
		Elements: []resp.Value{
			&resp.BulkString{Bytes: []byte("PSYNC")},
			&resp.BulkString{Bytes: []byte("?")},
			&resp.BulkString{Bytes: []byte("-1")},
		},
	}

	if err = session.encoder.Write(psyncInitialMsg); err != nil {
		return errors.Join(fmt.Errorf("error sending PSYNC to master %s from replica", replicaOf), err)
	}

	session.writer.Flush()

	psyncResponse, err := session.decoder.Read()
	if err != nil {
		return errors.Join(fmt.Errorf("error reading PSYNC response from master %s", replicaOf), err)
	}

	ss, ok := psyncResponse.(*resp.SimpleString)
	if !ok {
		return fmt.Errorf("ERR invalid response from master")
	}

	psyncParts := bytes.Split(ss.Bytes, []byte(" "))

	if len(psyncParts) < 3 || !slices.Equal(psyncParts[0], []byte("FULLRESYNC")) {
		return fmt.Errorf("ERR invalid response from master. Expected PSYNC response to contain at least 3 elements and FULLRESYNC")
	}

	s, err := session.reader.ReadString(byte('\n'))
	if err != nil {
		return fmt.Errorf("ERR invalid response from master. Expected rdb file after PSYNC")
	}

	s = strings.TrimLeft(s, "$")
	s = strings.TrimRight(s, "\r\n")

	rdbContentLen, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("ERR invalid RDB content len")
	}

	rdbContent := make([]byte, rdbContentLen)

	if _, err := io.ReadFull(session.reader, rdbContent); err != nil {
		return fmt.Errorf("ERR while reading RDB content")
	}

	session.Run()

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
		slog.String("address", l.Addr().String()),
		slog.Int("port", r.port),
	)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go r.acceptConnections()

	<-sigChan
	logger.Info("Shutdown signal received, stopping server...")

	// Stop accepting new connections
	err = l.Close()
	if err != nil {
		logger.Error("error closing listener", err.Error())
		return err
	}

	// Close all replica connections to unblock their session goroutines
	r.replicasRegistry.CloseAllConnections()

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

	session := NewSession(conn, r.store, r.transactions, r.isReplica, r.replicasRegistry, r.replicationId, false)
	session.Run()
}
