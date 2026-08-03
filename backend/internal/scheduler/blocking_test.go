package scheduler

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func TestSequentialSchedulerBlocksFailedBranch(t *testing.T) {
	calls := make([]string, 0, 2)
	executor := executorFunc(func(_ context.Context, definition pipeline.JobDefinition) error {
		calls = append(calls, definition.ID)
		if definition.ID == "failing" {
			return errors.New("permanent failure")
		}
		return nil
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	execution := newExecution(
		jobDefinition("grandchild", "child"),
		jobDefinition("child", "failing"),
		jobDefinition("failing"),
		jobDefinition("independent"),
	)

	if err := scheduler.Execute(context.Background(), &execution); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if execution.Status != runmodel.RunStatusFailed {
		t.Errorf("run status = %q, want failed", execution.Status)
	}
	wantCalls := []string{"failing", "independent"}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Errorf("executor calls = %v, want %v", calls, wantCalls)
	}

	wantStatuses := map[string]runmodel.JobStatus{
		"failing": runmodel.JobStatusFailed, "independent": runmodel.JobStatusCompleted,
		"child": runmodel.JobStatusBlocked, "grandchild": runmodel.JobStatusBlocked,
	}
	for _, job := range execution.Jobs {
		if job.Status != wantStatuses[job.JobID] {
			t.Errorf("job %q status = %q, want %q", job.JobID, job.Status, wantStatuses[job.JobID])
		}
	}
}
