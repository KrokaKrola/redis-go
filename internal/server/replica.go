package server

import (
	"fmt"
	"sync"

	"github.com/codecrafters-io/redis-starter-go/internal/replica"
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
