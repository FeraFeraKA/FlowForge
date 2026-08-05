package scheduler

import (
	"context"
	"sync"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestConcurrentSchedulerRespectsMaxParallel(t *testing.T) {
	var mu sync.Mutex
	active, maximum := 0, 0
	executor := executorFunc(func(ctx context.Context, _ pipeline.JobDefinition) error {
		mu.Lock()
		active++
		if active > maximum {
			maximum = active
		}
		mu.Unlock()
		defer func() {
			mu.Lock()
			active--
			mu.Unlock()
		}()

		timer := time.NewTimer(20 * time.Millisecond)
		select {
		case <-timer.C:
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		}
		return nil
	})
	scheduler := NewConcurrentScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	}, nil)
	execution := newExecution(
		jobDefinition("one"), jobDefinition("two"),
		jobDefinition("three"), jobDefinition("four"),
	)
	execution.MaxParallel = 2

	if err := scheduler.Execute(context.Background(), &execution); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if maximum != 2 {
		t.Errorf("maximum parallel executions = %d, want 2", maximum)
	}
	if execution.Status != runmodel.RunStatusCompleted {
		t.Errorf("run status = %q, want completed", execution.Status)
	}
	for _, job := range execution.Jobs {
		if job.Status != runmodel.JobStatusCompleted || job.Attempt != 1 {
			t.Errorf("job %q = (%q, attempt %d)", job.JobID, job.Status, job.Attempt)
		}
	}
}
