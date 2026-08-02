package pipeline

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestMemoryRepositoryDoesNotExposeStoredData(t *testing.T) {
	repository := NewMemoryRepository()
	pipeline := validPipeline()
	pipeline.Jobs[0].DependsOn = []string{"source"}
	mustCreatePipeline(t, repository, pipeline)

	pipeline.Name = "Changed input"
	pipeline.Jobs[0].Name = "Changed input job"
	pipeline.Jobs[0].DependsOn[0] = "changed-input"
	pipeline.Jobs[0].DelayConfig.Duration = 2 * time.Second

	first, err := repository.GetByID(context.Background(), pipeline.ID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}
	if first.Name == pipeline.Name || first.Jobs[0].Name == pipeline.Jobs[0].Name {
		t.Fatal("stored pipeline changed after mutating Create argument")
	}

	first.Name = "Changed result"
	first.Jobs[0].DelayConfig.Duration = 3 * time.Second
	second, err := repository.GetByID(context.Background(), pipeline.ID)
	if err != nil {
		t.Fatalf("second GetByID() error = %v", err)
	}
	if second.Name == first.Name || second.Jobs[0].DelayConfig.Duration == first.Jobs[0].DelayConfig.Duration {
		t.Error("stored pipeline changed after mutating GetByID result")
	}
}

func TestMemoryRepositoryListNewestFirst(t *testing.T) {
	repository := NewMemoryRepository()
	older := validPipeline()
	older.ID, older.CreatedAt = "older", time.Date(2026, time.August, 1, 12, 0, 0, 0, time.UTC)
	newer := validPipeline()
	newer.ID, newer.CreatedAt = "newer", older.CreatedAt.Add(time.Minute)
	mustCreatePipeline(t, repository, older)
	mustCreatePipeline(t, repository, newer)

	pipelines, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(pipelines) != 2 {
		t.Fatalf("List() length = %d, want 2", len(pipelines))
	}
	if pipelines[0].ID != "newer" || pipelines[1].ID != "older" {
		t.Errorf("List() IDs = %v, want [newer older]", []string{pipelines[0].ID, pipelines[1].ID})
	}
}

func mustCreatePipeline(t *testing.T, repository *MemoryRepository, pipeline Pipeline) {
	t.Helper()
	if err := repository.Create(context.Background(), pipeline); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
}

func TestMemoryRepositoryConcurrentAccess(t *testing.T) {
	const pipelineCount = 50
	repository := NewMemoryRepository()
	errCh := make(chan error, pipelineCount*2)
	var workers sync.WaitGroup

	for index := range pipelineCount {
		workers.Add(2)
		go func() {
			defer workers.Done()
			pipeline := validPipeline()
			pipeline.ID = fmt.Sprintf("pipeline-%d", index)
			if err := repository.Create(context.Background(), pipeline); err != nil {
				errCh <- err
			}
		}()
		go func() {
			defer workers.Done()
			if _, err := repository.List(context.Background()); err != nil {
				errCh <- err
			}
		}()
	}

	workers.Wait()
	close(errCh)
	for err := range errCh {
		t.Errorf("concurrent repository operation: %v", err)
	}

	pipelines, err := repository.List(context.Background())
	if err != nil {
		t.Fatalf("final List() error = %v", err)
	}
	if len(pipelines) != pipelineCount {
		t.Errorf("final pipeline count = %d, want %d", len(pipelines), pipelineCount)
	}
}
