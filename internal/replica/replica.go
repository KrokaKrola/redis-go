package replica

import (
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

type Replica struct {
	Port         int
	Capabilities []string
	Connection   net.Conn
	AckOffset    int // last acknowledged offset from this replica
}

type ReplicasRegistry interface {
	AddReplica(addr string, port int) error
	AddCapabilities(addr string, capabilities []string) error
	GetReplica(addr string) (*Replica, bool)
	AddReplicaConnection(conn net.Conn) error
	BroadcastRespValue(value resp.Value)
	GetAllReplicas() []*Replica
	UpdateAckOffset(addr string, offset int) error
}
