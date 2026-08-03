package scheduler

import (
	"context"
	"errors"
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

func TestDelayExecutorReturnsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := (DelayExecutor{}).Execute(ctx, delayDefinition(time.Second))
	if !errors.Is(err, context.Canceled) {
		t.Errorf("Execute() error = %v, want context.Canceled", err)
	}
}

func TestDelayExecutorReturnsDeadline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	err := (DelayExecutor{}).Execute(ctx, delayDefinition(time.Second))
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Execute() error = %v, want context.DeadlineExceeded", err)
	}
}

func TestDelayExecutorRejectsInvalidConfig(t *testing.T) {
	tests := []struct {
		name       string
		definition pipeline.JobDefinition
	}{
		{
			name: "wrong executor type",
			definition: pipeline.JobDefinition{
				Type: "http", DelayConfig: &pipeline.DelayConfig{Duration: time.Second},
			},
		},
		{
			name:       "missing config",
			definition: pipeline.JobDefinition{Type: pipeline.ExecutorTypeDelay},
		},
		{name: "zero duration", definition: delayDefinition(0)},
		{name: "negative duration", definition: delayDefinition(-time.Second)},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := (DelayExecutor{}).Execute(context.Background(), test.definition)
			if !errors.Is(err, ErrInvalidExecutorConfig) {
				t.Errorf("Execute() error = %v, want ErrInvalidExecutorConfig", err)
			}
		})
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
