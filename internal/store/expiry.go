package store

import "strings"

type ExpiryType string

const (
	EXPIRY_PX ExpiryType = "px"
	EXPIRY_EX ExpiryType = "ex"
)

func ProcessExpType(v string) (ExpiryType, bool) {
	vLower := strings.ToLower(v)

	// allow empty string for indefinite key storage
	if vLower == "" {
		return "", true
	}

	if vLower != string(EXPIRY_EX) && vLower != string(EXPIRY_PX) {
		return "", false
	}

	return ExpiryType(vLower), true
}
