package scheduler

import (
	"context"
	"errors"
	"testing"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestSequentialSchedulerCancelsRunDuringExecution(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	executor := executorFunc(func(executionCtx context.Context, _ pipeline.JobDefinition) error {
		cancel()
		<-executionCtx.Done()
		return executionCtx.Err()
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	execution := newExecution(
		jobDefinition("first"),
		jobDefinition("second", "first"),
	)

	err := scheduler.Execute(ctx, &execution)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("Execute() error = %v, want context.Canceled", err)
	}
	if execution.Status != runmodel.RunStatusCancelled {
		t.Errorf("run status = %q, want cancelled", execution.Status)
	}
	for _, job := range execution.Jobs {
		if job.Status != runmodel.JobStatusCancelled {
			t.Errorf("job %q status = %q, want cancelled", job.JobID, job.Status)
		}
		if job.FinishedAt == nil {
			t.Errorf("job %q FinishedAt = nil", job.JobID)
		}
	}
	if execution.StartedAt == nil || execution.FinishedAt == nil {
		t.Error("run lifecycle timestamps were not set")
	}
}
