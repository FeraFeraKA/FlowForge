package run

import (
	"context"
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
)

func TestMemoryRepositoryOwnsRunSnapshots(t *testing.T) {
	repository := NewMemoryRepository()
	execution := sampleRun("run-1", time.Now())
	if err := repository.Create(context.Background(), execution); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	execution.Jobs[0].Status = JobStatusFailed

	stored, err := repository.GetByID(context.Background(), execution.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if stored.Jobs[0].Status != JobStatusPending {
		t.Error("Create() stored shared JobRun data")
	}

	stored.Jobs[0].Status = JobStatusCompleted
	if err := repository.Save(context.Background(), stored); err != nil {
		t.Fatalf("Save() error = %v", err)
	}
	stored.Jobs[0].Status = JobStatusFailed
	updated, err := repository.GetByID(context.Background(), execution.ID)
	if err != nil {
		t.Fatalf("second GetByID() error = %v", err)
	}
	if updated.Jobs[0].Status != JobStatusCompleted {
		t.Error("Save() stored shared JobRun data")
	}
}

func TestMemoryRepositoryListsNewestFirst(t *testing.T) {
	repository := NewMemoryRepository()
	older := sampleRun("older", time.Now())
	newer := sampleRun("newer", older.CreatedAt.Add(time.Second))
	if err := repository.Create(context.Background(), older); err != nil {
		t.Fatalf("Create(older) error = %v", err)
	}
	if err := repository.Create(context.Background(), newer); err != nil {
		t.Fatalf("Create(newer) error = %v", err)
	}

	runs, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(runs) != 2 {
		t.Fatalf("List() length = %d, want 2", len(runs))
	}
	if runs[0].ID != "newer" || runs[1].ID != "older" {
		t.Errorf("List() = %#v, want [newer older]", runs)
	}
}

func sampleRun(id string, createdAt time.Time) Run {
	return New(id, pipeline.Pipeline{
		ID: "pipeline", Name: "Pipeline", MaxParallel: 1,
		Jobs: []pipeline.JobDefinition{{ID: "job"}},
	}, createdAt)
}
