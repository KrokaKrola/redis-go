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

func getPossibleEndTime() time.Time {
	return time.Date(9999, 12, 31, 23, 59, 59, 999, time.UTC)
}
