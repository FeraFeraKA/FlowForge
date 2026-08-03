package scheduler

import (
	"context"
	"errors"
	"testing"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestSequentialSchedulerRejectsInvalidSnapshot(t *testing.T) {
	mismatched := newExecution(jobDefinition("job"))
	mismatched.Jobs = nil
	duplicateDefinitions := newExecution(jobDefinition("same"), jobDefinition("same"))
	unknownState := newExecution(jobDefinition("known"))
	unknownState.Jobs[0].JobID = "unknown"
	empty := runmodel.Run{}

	tests := []struct {
		name      string
		execution *runmodel.Run
	}{
		{name: "nil run", execution: nil},
		{name: "empty run", execution: &empty},
		{name: "mismatched counts", execution: &mismatched},
		{name: "duplicate definitions", execution: &duplicateDefinitions},
		{name: "unknown job state", execution: &unknownState},
	}
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executorFunc(func(context.Context, pipeline.JobDefinition) error {
			return nil
		}),
	})

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := scheduler.Execute(context.Background(), test.execution)
			if !errors.Is(err, ErrInvalidRun) {
				t.Errorf("Execute() error = %v, want ErrInvalidRun", err)
			}
		})
	}
}

func TestSequentialSchedulerRequiresRegisteredExecutor(t *testing.T) {
	scheduler := NewSequentialScheduler(nil)
	execution := newExecution(jobDefinition("job"))

	err := scheduler.Execute(context.Background(), &execution)
	if !errors.Is(err, ErrExecutorNotRegistered) {
		t.Errorf("Execute() error = %v, want ErrExecutorNotRegistered", err)
	}
}

func TestSequentialSchedulerReportsStalledGraph(t *testing.T) {
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executorFunc(func(context.Context, pipeline.JobDefinition) error {
			return nil
		}),
	})
	execution := newExecution(
		jobDefinition("first", "second"),
		jobDefinition("second", "first"),
	)

	err := scheduler.Execute(context.Background(), &execution)
	if !errors.Is(err, ErrSchedulerStalled) {
		t.Errorf("Execute() error = %v, want ErrSchedulerStalled", err)
	}
}
