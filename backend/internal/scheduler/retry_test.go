package scheduler

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestSequentialSchedulerRetriesFailedJob(t *testing.T) {
	attempts := 0
	executor := executorFunc(func(_ context.Context, _ pipeline.JobDefinition) error {
		attempts++
		if attempts == 1 {
			return errors.New("temporary failure")
		}
		return nil
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	definition := jobDefinition("unstable")
	definition.MaxRetries = 1
	execution := newExecution(definition)

	if err := scheduler.Execute(context.Background(), &execution); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	job := execution.Jobs[0]
	if execution.Status != runmodel.RunStatusCompleted || job.Status != runmodel.JobStatusCompleted {
		t.Errorf("statuses = (%q, %q), want (completed, completed)", execution.Status, job.Status)
	}
	if job.Attempt != 2 || attempts != 2 || job.Error != "" {
		t.Errorf("retry result = attempt %d, calls %d, error %q", job.Attempt, attempts, job.Error)
	}
}

func TestSequentialSchedulerRetriesTimedOutAttempt(t *testing.T) {
	executor := executorFunc(func(ctx context.Context, _ pipeline.JobDefinition) error {
		<-ctx.Done()
		return ctx.Err()
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	definition := jobDefinition("timeout")
	definition.Timeout = 5 * time.Millisecond
	definition.MaxRetries = 1
	execution := newExecution(definition)

	if err := scheduler.Execute(context.Background(), &execution); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	job := execution.Jobs[0]
	if execution.Status != runmodel.RunStatusFailed || job.Status != runmodel.JobStatusFailed {
		t.Errorf("statuses = (%q, %q), want (failed, failed)", execution.Status, job.Status)
	}
	if job.Attempt != 2 || !strings.Contains(job.Error, context.DeadlineExceeded.Error()) {
		t.Errorf("timeout result = attempt %d, error %q", job.Attempt, job.Error)
	}
}
