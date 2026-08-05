package scheduler

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestConcurrentSchedulerRetriesWithoutBlockingOtherJobs(t *testing.T) {
	var unstableAttempts atomic.Int32
	executor := executorFunc(func(_ context.Context, definition pipeline.JobDefinition) error {
		if definition.ID == "unstable" && unstableAttempts.Add(1) == 1 {
			return errors.New("temporary failure")
		}
		return nil
	})
	scheduler := NewConcurrentScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	}, nil)
	unstable := jobDefinition("unstable")
	unstable.MaxRetries = 1
	unstable.RetryDelay = time.Millisecond
	execution := newExecution(unstable, jobDefinition("independent"))
	execution.MaxParallel = 2

	if err := scheduler.Execute(context.Background(), &execution); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if execution.Status != runmodel.RunStatusCompleted {
		t.Errorf("run status = %q, want completed", execution.Status)
	}
	if execution.Jobs[0].Status != runmodel.JobStatusCompleted || execution.Jobs[0].Attempt != 2 {
		t.Errorf("unstable job = (%q, attempt %d)", execution.Jobs[0].Status, execution.Jobs[0].Attempt)
	}
	if execution.Jobs[1].Status != runmodel.JobStatusCompleted {
		t.Errorf("independent job status = %q, want completed", execution.Jobs[1].Status)
	}
}
