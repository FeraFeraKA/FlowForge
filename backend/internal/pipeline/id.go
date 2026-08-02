package pipeline

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func generatePipelineID() (string, error) {
	var randomBytes [16]byte
	if _, err := rand.Read(randomBytes[:]); err != nil {
		return "", fmt.Errorf("generate random pipeline ID: %w", err)
	}

	return "pl_" + hex.EncodeToString(randomBytes[:]), nil
}
