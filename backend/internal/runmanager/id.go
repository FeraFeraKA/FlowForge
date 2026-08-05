package runmanager

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generateRunID() (string, error) {
	var randomBytes [16]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		return "", fmt.Errorf("generate random run ID: %w", err)
	}
	return "run_" + hex.EncodeToString(randomBytes[:]), nil
}
