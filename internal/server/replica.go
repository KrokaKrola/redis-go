package server

import (
	"fmt"
	"sync"
)

type replica struct {
	port         int
	capabilities []string
}

type ReplicasRegistry struct {
	sync.RWMutex
	registry map[string]*replica
}

func NewReplicasRegistry() *ReplicasRegistry {
	return &ReplicasRegistry{
		registry: make(map[string]*replica),
	}
}

func (rr *ReplicasRegistry) AddReplica(address string, port int) error {
	rr.Lock()
	defer rr.Unlock()

	_, ok := rr.registry[address]
	if ok {
		return fmt.Errorf("ERR replica has already been assigned to this port")
	}

	rr.registry[address] = &replica{
		port: port,
	}

	return nil
}

func (rr *ReplicasRegistry) AddCapabilities(address string, capabilities []string) error {
	rr.Lock()
	defer rr.Unlock()

	replica, ok := rr.registry[address]
	if !ok {
		return fmt.Errorf("ERR replica is not found")
	}

	newCapabilities := append(replica.capabilities, capabilities...)
	rr.registry[address].capabilities = newCapabilities

	return nil
}
