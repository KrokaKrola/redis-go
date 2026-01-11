package replica

type ReplicasRegistry interface {
	AddReplica(addr string, port int) error
	AddCapabilities(addr string, capabilities []string) error
}
