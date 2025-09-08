package server

import (
	"fmt"
	"net"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

func processRespCommand(conn net.Conn, r *resp.Resp) error {
	switch r.Command {
	case resp.ECHO_COMMAND:
		logger.Debug("Processing ECHO command")
		responseCmd := resp.NewResp(resp.ECHO_COMMAND, r.Value)

		s, err := resp.ToString(responseCmd)
		if err != nil {
			return err
		}

		conn.Write([]byte(s))

		return nil
	default:
		return fmt.Errorf("unknown resp command: %s\n", string(r.Command))
	}
}
