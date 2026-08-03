package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestSequentialSchedulerCancelsBeforeStart(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	executorCalled := false
	executor := executorFunc(func(context.Context, pipeline.JobDefinition) error {
		executorCalled = true
		return nil
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	execution := newExecution(jobDefinition("pending"))

	err := scheduler.Execute(ctx, &execution)
	if !errors.Is(err, context.Canceled) || execution.Status != runmodel.RunStatusCancelled {
		t.Fatalf("Execute() = (%q, %v), want (cancelled, context.Canceled)", execution.Status, err)
	}
	if executorCalled || execution.StartedAt != nil {
		t.Error("cancelled-before-start run started execution")
	}
}

func TestSequentialSchedulerCancelsRetryDelay(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	executor := executorFunc(func(context.Context, pipeline.JobDefinition) error {
		return errors.New("retryable failure")
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	definition := jobDefinition("retrying")
	definition.MaxRetries = 1
	definition.RetryDelay = time.Second
	execution := newExecution(definition)

	err := scheduler.Execute(ctx, &execution)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Execute() error = %v, want context.DeadlineExceeded", err)
	}
	if execution.Status != runmodel.RunStatusCancelled ||
		execution.Jobs[0].Status != runmodel.JobStatusCancelled {
		t.Errorf("statuses = (%q, %q), want cancelled", execution.Status, execution.Jobs[0].Status)
	}
	if execution.Jobs[0].Attempt != 1 {
		t.Errorf("attempt = %d, want 1", execution.Jobs[0].Attempt)
	}
}
