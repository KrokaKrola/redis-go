package commands

import (
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
	"github.com/codecrafters-io/redis-starter-go/internal/store"
)

// testDispatch is a helper for tests that wraps the new Dispatch signature.
// It creates contexts with sensible defaults for testing.
func testDispatch(cmd *Command, s *store.Store, isReplica bool) resp.Value {
	return Dispatch(
		&ServerContext{
			Store:     s,
			IsReplica: isReplica,
		},
		&HandlerContext{
			Cmd: cmd,
		},
	)
}
