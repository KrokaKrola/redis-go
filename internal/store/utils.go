package store

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"
)

func newId() string {
	t := time.Now().UnixNano()
	b := make([]byte, 4)
	rand.Read(b)
	return fmt.Sprintf("%x-%s", t, hex.EncodeToString(b))
}
