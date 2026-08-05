package scheduler

import (
	"context"
	"errors"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestConcurrentSchedulerCancelsAndWaitsForWorkers(t *testing.T) {
	started := make(chan struct{}, 2)
	executor := executorFunc(func(ctx context.Context, _ pipeline.JobDefinition) error {
		started <- struct{}{}
		<-ctx.Done()
		return ctx.Err()
	})
	scheduler := NewConcurrentScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	}, nil)
	execution := newExecution(
		jobDefinition("one"), jobDefinition("two"), jobDefinition("three"),
	)
	execution.MaxParallel = 2
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		done <- scheduler.Execute(ctx, &execution)
	}()

	<-started
	<-started
	cancel()
	timer := time.NewTimer(time.Second)
	defer timer.Stop()

	select {
	case err := <-done:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Execute() error = %v, want context.Canceled", err)
		}
	case <-timer.C:
		t.Fatal("Execute() did not wait for cancelled workers")
	}
	if execution.Status != runmodel.RunStatusCancelled {
		t.Errorf("run status = %q, want cancelled", execution.Status)
	}
	for _, job := range execution.Jobs {
		if job.Status != runmodel.JobStatusCancelled {
			t.Errorf("job %q status = %q, want cancelled", job.JobID, job.Status)
		}
	}
}
