package replica

type Replica struct {
	Port         int
	Capabilities []string
}

type ReplicasRegistry interface {
	AddReplica(addr string, port int) error
	AddCapabilities(addr string, capabilities []string) error
	GetReplica(addr string) (*Replica, bool)
}
