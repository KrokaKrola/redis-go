package commands

type Name string

const (
	PING_COMMAND  Name = "PING"
	ECHO_COMMAND  Name = "ECHO"
	GET_COMMAND   Name = "GET"
	SET_COMMAND   Name = "SET"
	RPUSH_COMMAND Name = "RPUSH"
)
