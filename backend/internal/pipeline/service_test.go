package pipeline

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestServiceCreatesPipeline(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewService(repository)
	createdAt := time.Date(2026, time.August, 2, 10, 0, 0, 0, time.UTC)
	service.newID = func() (string, error) { return "pl_test", nil }
	service.now = func() time.Time { return createdAt }
	template := validPipeline()

	created, err := service.Create(context.Background(), CreateParams{
		Name: template.Name, MaxParallel: template.MaxParallel, Jobs: template.Jobs,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if created.ID != "pl_test" || !created.CreatedAt.Equal(createdAt) {
		t.Errorf("created identity = (%q, %v), want (pl_test, %v)", created.ID, created.CreatedAt, createdAt)
	}

	template.Jobs[0].Name = "Changed input"
	stored, err := service.GetByID(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if stored.Jobs[0].Name == template.Jobs[0].Name {
		t.Error("created pipeline shares Jobs with CreateParams")
	}
}

func TestServiceRejectsInvalidPipelineBeforeSaving(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewService(repository)

	_, err := service.Create(context.Background(), CreateParams{})
	requireValidationErrors(t, err)
	pipelines, listErr := repository.List(context.Background())
	if listErr != nil {
		t.Fatalf("List() error = %v", listErr)
	}
	if len(pipelines) != 0 {
		t.Errorf("repository contains %d pipelines after rejected Create, want 0", len(pipelines))
	}
}

func TestServicePreservesRepositoryError(t *testing.T) {
	service := NewService(NewMemoryRepository())

	_, err := service.GetByID(context.Background(), "missing")
	if !errors.Is(err, ErrPipelineNotFound) {
		t.Errorf("GetByID() error = %v, want ErrPipelineNotFound", err)
	}
}
