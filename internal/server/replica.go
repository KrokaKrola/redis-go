package server

import (
	"bufio"
	"fmt"
	"log/slog"
	"net"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/replica"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type ReplicasRegistry struct {
	sync.RWMutex
	registry map[string]*replica.Replica
}

func NewReplicasRegistry() *ReplicasRegistry {
	return &ReplicasRegistry{
		registry: make(map[string]*replica.Replica),
	}
}

func (rr *ReplicasRegistry) AddReplica(address string, port int) error {
	rr.Lock()
	defer rr.Unlock()

	_, ok := rr.registry[address]
	if ok {
		return fmt.Errorf("ERR replica has already been assigned to this port")
	}

	rr.registry[address] = &replica.Replica{
		Port: port,
	}

	return nil
}

func (rr *ReplicasRegistry) AddCapabilities(address string, capabilities []string) error {
	rr.Lock()
	defer rr.Unlock()

	replItem, ok := rr.registry[address]
	if !ok {
		return fmt.Errorf("ERR replica is not found")
	}

	rr.registry[address].Capabilities = append(replItem.Capabilities, capabilities...)

	return nil
}

func (rr *ReplicasRegistry) GetReplica(addr string) (*replica.Replica, bool) {
	rr.RLock()
	defer rr.RUnlock()

	replItem, ok := rr.registry[addr]

	return replItem, ok
}

func (rr *ReplicasRegistry) AddReplicaConnection(conn net.Conn) error {
	addr := conn.RemoteAddr().String()
	replItem, ok := rr.registry[addr]

	if !ok {
		return fmt.Errorf("ERR replica is not found")
	}

	replItem.Connection = conn

	return nil
}

func (rr *ReplicasRegistry) BroadcastRespValue(value resp.Value) {
	for _, v := range rr.registry {
		writer := bufio.NewWriter(v.Connection)
		encoder := resp.NewEncoder(writer)

		logger.Debug("sending message to replica", slog.String("address", v.Connection.RemoteAddr().String()))

		if err := encoder.Write(value); err != nil {
			logger.Error("error while trying to broadcast resp value to", "address", v.Connection.RemoteAddr().String())
			continue
		}
		writer.Flush()
	}
}
