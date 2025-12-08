package commands

import "strings"

type Name string

const (
	PING_COMMAND    Name = "PING"
	ECHO_COMMAND    Name = "ECHO"
	GET_COMMAND     Name = "GET"
	SET_COMMAND     Name = "SET"
	RPUSH_COMMAND   Name = "RPUSH"
	LRANGE_COMMAND  Name = "LRANGE"
	LPUSH_COMMAND   Name = "LPUSH"
	LLEN_COMMAND    Name = "LLEN"
	LPOP_COMMAND    Name = "LPOP"
	BLPOP_COMMAND   Name = "BLPOP"
	TYPE_COMMAND    Name = "TYPE"
	XADD_COMMAND    Name = "XADD"
	XRANGE_COMMAND  Name = "XRANGE"
	XREAD_COMMAND   Name = "XREAD"
	INCR_COMMAND    Name = "INCR"
	MULTI_COMMAND   Name = "MULTI"
	EXEC_COMMAND    Name = "EXEC"
	DISCARD_COMMAND Name = "DISCARD"
	INFO_COMMAND    Name = "INFO"
)

var commandByName = map[string]Name{
	string(PING_COMMAND):    PING_COMMAND,
	string(ECHO_COMMAND):    ECHO_COMMAND,
	string(GET_COMMAND):     GET_COMMAND,
	string(SET_COMMAND):     SET_COMMAND,
	string(RPUSH_COMMAND):   RPUSH_COMMAND,
	string(LRANGE_COMMAND):  LRANGE_COMMAND,
	string(LPUSH_COMMAND):   LPUSH_COMMAND,
	string(LLEN_COMMAND):    LLEN_COMMAND,
	string(LPOP_COMMAND):    LPOP_COMMAND,
	string(BLPOP_COMMAND):   BLPOP_COMMAND,
	string(TYPE_COMMAND):    TYPE_COMMAND,
	string(XADD_COMMAND):    XADD_COMMAND,
	string(XRANGE_COMMAND):  XRANGE_COMMAND,
	string(XREAD_COMMAND):   XREAD_COMMAND,
	string(INCR_COMMAND):    INCR_COMMAND,
	string(MULTI_COMMAND):   MULTI_COMMAND,
	string(EXEC_COMMAND):    EXEC_COMMAND,
	string(DISCARD_COMMAND): DISCARD_COMMAND,
	string(INFO_COMMAND):    INFO_COMMAND,
}

func getCommandName(name []byte) Name {
	upper := strings.ToUpper(string(name))
	return commandByName[upper]
}
