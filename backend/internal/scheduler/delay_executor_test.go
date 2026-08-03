package scheduler

import (
	"context"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
)

func TestDelayExecutorWaitsForDuration(t *testing.T) {
	const duration = 20 * time.Millisecond
	startedAt := time.Now()

	err := (DelayExecutor{}).Execute(context.Background(), delayDefinition(duration))
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if elapsed := time.Since(startedAt); elapsed < duration {
		t.Errorf("Execute() elapsed time = %v, want at least %v", elapsed, duration)
	}
}

func delayDefinition(duration time.Duration) pipeline.JobDefinition {
	return pipeline.JobDefinition{
		ID:          "delay",
		Name:        "Delay",
		Type:        pipeline.ExecutorTypeDelay,
		DelayConfig: &pipeline.DelayConfig{Duration: duration},
	}
}
