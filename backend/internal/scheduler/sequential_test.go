package scheduler

import (
	"context"
	"reflect"
	"testing"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

type executorFunc func(context.Context, pipeline.JobDefinition) error

func (function executorFunc) Execute(ctx context.Context, definition pipeline.JobDefinition) error {
	return function(ctx, definition)
}

func TestSequentialSchedulerCompletesDiamondDAG(t *testing.T) {
	calls := make([]string, 0, 4)
	executor := executorFunc(func(_ context.Context, definition pipeline.JobDefinition) error {
		calls = append(calls, definition.ID)
		return nil
	})
	scheduler := NewSequentialScheduler(map[pipeline.ExecutorType]Executor{
		pipeline.ExecutorTypeDelay: executor,
	})
	execution := newExecution(
		jobDefinition("install"),
		jobDefinition("tests", "install"),
		jobDefinition("build", "install"),
		jobDefinition("deploy", "tests", "build"),
	)

	if err := scheduler.Execute(context.Background(), &execution); err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if execution.Status != runmodel.RunStatusCompleted {
		t.Errorf("run status = %q, want completed", execution.Status)
	}
	wantCalls := []string{"install", "tests", "build", "deploy"}
	if !reflect.DeepEqual(calls, wantCalls) {
		t.Errorf("executor calls = %v, want %v", calls, wantCalls)
	}
	for _, job := range execution.Jobs {
		if job.Status != runmodel.JobStatusCompleted || job.Attempt != 1 || job.Progress != 100 {
			t.Errorf("job %q state = (%q, attempt %d, progress %d)", job.JobID, job.Status, job.Attempt, job.Progress)
		}
	}
	if execution.StartedAt == nil || execution.FinishedAt == nil {
		t.Error("run timestamps were not set")
	}
}

func newExecution(definitions ...pipeline.JobDefinition) runmodel.Run {
	jobs := make([]runmodel.JobRun, len(definitions))
	for index, definition := range definitions {
		jobs[index] = runmodel.JobRun{JobID: definition.ID, Status: runmodel.JobStatusPending}
	}
	return runmodel.Run{
		ID: "run-1", Status: runmodel.RunStatusPending,
		JobDefinitions: definitions, Jobs: jobs,
	}
}

func jobDefinition(id string, dependencies ...string) pipeline.JobDefinition {
	return pipeline.JobDefinition{
		ID: id, Name: id, Type: pipeline.ExecutorTypeDelay, DependsOn: dependencies,
	}
}
