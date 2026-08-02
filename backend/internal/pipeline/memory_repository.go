package pipeline

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

var (
	ErrPipelineNotFound      = errors.New("pipeline not found")
	ErrPipelineAlreadyExists = errors.New("pipeline already exists")
)

type MemoryRepository struct {
	mu        sync.RWMutex
	pipelines map[string]Pipeline
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{pipelines: make(map[string]Pipeline)}
}

func (r *MemoryRepository) Create(ctx context.Context, pipeline Pipeline) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	stored := clonePipeline(pipeline)
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.pipelines[pipeline.ID]; exists {
		return fmt.Errorf("%w: %s", ErrPipelineAlreadyExists, pipeline.ID)
	}

	r.pipelines[pipeline.ID] = stored
	return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (Pipeline, error) {
	if err := ctx.Err(); err != nil {
		return Pipeline{}, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	stored, exists := r.pipelines[id]
	if !exists {
		return Pipeline{}, fmt.Errorf("%w: %s", ErrPipelineNotFound, id)
	}

	return clonePipeline(stored), nil
}

func (r *MemoryRepository) List(ctx context.Context) ([]Pipeline, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	r.mu.RLock()
	pipelines := make([]Pipeline, 0, len(r.pipelines))
	for _, stored := range r.pipelines {
		pipelines = append(pipelines, clonePipeline(stored))
	}
	r.mu.RUnlock()

	sort.Slice(pipelines, func(left, right int) bool {
		if pipelines[left].CreatedAt.Equal(pipelines[right].CreatedAt) {
			return pipelines[left].ID < pipelines[right].ID
		}
		return pipelines[left].CreatedAt.After(pipelines[right].CreatedAt)
	})

	return pipelines, nil
}
