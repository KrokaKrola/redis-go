package main

import (
	"log/slog"
	"os"

	"github.com/codecrafters-io/redis-starter-go/internal/logger"
	"github.com/codecrafters-io/redis-starter-go/internal/server"
)

const port uint16 = 6379

func main() {
	logger.Init(slog.LevelDebug, "text")

	s := server.NewRedisServer(port)

	err := s.Listen()

	if err != nil {
		logger.Error("error starting the server", "err", err)
		os.Exit(1)
	}
}
