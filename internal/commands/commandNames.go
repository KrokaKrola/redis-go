package commands

import "strings"

type Name string

const (
	PING_COMMAND   Name = "PING"
	ECHO_COMMAND   Name = "ECHO"
	GET_COMMAND    Name = "GET"
	SET_COMMAND    Name = "SET"
	RPUSH_COMMAND  Name = "RPUSH"
	LRANGE_COMMAND Name = "LRANGE"
	LPUSH_COMMAND  Name = "LPUSH"
	LLEN_COMMAND   Name = "LLEN"
)

var commandByName = map[string]Name{
	string(PING_COMMAND):   PING_COMMAND,
	string(ECHO_COMMAND):   ECHO_COMMAND,
	string(GET_COMMAND):    GET_COMMAND,
	string(SET_COMMAND):    SET_COMMAND,
	string(RPUSH_COMMAND):  RPUSH_COMMAND,
	string(LRANGE_COMMAND): LRANGE_COMMAND,
	string(LPUSH_COMMAND):  LPUSH_COMMAND,
	string(LLEN_COMMAND):   LLEN_COMMAND,
}

func getCommandName(name []byte) Name {
	upper := strings.ToUpper(string(name))
	return commandByName[upper]
}
