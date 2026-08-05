package runmanager

import (
	"context"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
	"flowforge/backend/internal/scheduler"
)

func TestManagerRunsPipelineInBackground(t *testing.T) {
	manager, runs := testManager(t, 5*time.Millisecond)

	initial, err := manager.StartRun(context.Background(), "pipeline")
	if err != nil {
		t.Fatalf("StartRun() error = %v", err)
	}
	if initial.Status != runmodel.RunStatusPending {
		t.Errorf("initial status = %q, want pending", initial.Status)
	}
	manager.Wait()

	finished, err := runs.GetByID(context.Background(), initial.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if finished.Status != runmodel.RunStatusCompleted ||
		finished.Jobs[0].Status != runmodel.JobStatusCompleted {
		t.Errorf("finished statuses = (%q, %q)", finished.Status, finished.Jobs[0].Status)
	}
}

func testManager(t *testing.T, duration time.Duration) (*Manager, *runmodel.MemoryRepository) {
	t.Helper()
	pipelines := pipeline.NewMemoryRepository()
	definition := pipeline.Pipeline{
		ID: "pipeline", Name: "Pipeline", MaxParallel: 1,
		Jobs: []pipeline.JobDefinition{{
			ID: "job", Name: "Job", Type: pipeline.ExecutorTypeDelay,
			DelayConfig: &pipeline.DelayConfig{Duration: duration},
		}},
	}
	if err := pipelines.Create(context.Background(), definition); err != nil {
		t.Fatalf("create pipeline: %v", err)
	}
	runs := runmodel.NewMemoryRepository()
	concurrent := scheduler.NewConcurrentScheduler(
		map[pipeline.ExecutorType]scheduler.Executor{
			pipeline.ExecutorTypeDelay: scheduler.DelayExecutor{},
		},
		func(snapshot runmodel.Run) error {
			return runs.Save(context.Background(), snapshot)
		},
	)
	manager := New(context.Background(), pipelines, runs, concurrent)
	manager.newID = func() (string, error) { return "run", nil }
	return manager, runs
}
