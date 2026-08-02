package pipeline

import (
	"strings"
	"testing"
)

func TestGeneratePipelineID(t *testing.T) {
	seen := make(map[string]struct{}, 100)
	for range 100 {
		id, err := generatePipelineID()
		if err != nil {
			t.Fatalf("generatePipelineID() error = %v", err)
		}
		if !strings.HasPrefix(id, "pl_") || len(id) != 35 {
			t.Fatalf("generatePipelineID() = %q, want pl_ followed by 32 hex characters", id)
		}
		if _, duplicate := seen[id]; duplicate {
			t.Fatalf("generatePipelineID() returned duplicate %q", id)
		}
		seen[id] = struct{}{}
	}
}
