package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/server"
)

const defaultPortValue int = 6379

func main() {
	logger.Init(slog.LevelDebug, "text")

	var port int
	var replicaOf string

	flag.IntVar(&port, "port", defaultPortValue, "Defines port number for redis server")
	flag.StringVar(&replicaOf, "replicaof", "", "Defines replica host and port")
	flag.Parse()

	if port < 1 || port > 65535 {
		logger.Error("invalid port number", "port", port)
		os.Exit(1)
	}

	isReplica := replicaOf != ""

	s := server.NewRedisServer(port, isReplica)

	if isReplica {
		go func() {
			if err := s.ConnectToMaster(replicaOf, port); err != nil {
				logger.Error(err.Error(), "replicaOf", replicaOf, "port", port)
			}
		}()
	}

	err := s.Listen()

	if err != nil {
		logger.Error("error starting the server", "err", err)
		os.Exit(1)
	}
}
