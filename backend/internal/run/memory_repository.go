package run

import (
	"context"
	"errors"
	"fmt"
	"sync"
)

var (
	ErrRunNotFound      = errors.New("run not found")
	ErrRunAlreadyExists = errors.New("run already exists")
)

type MemoryRepository struct {
	mu   sync.RWMutex
	runs map[string]Run
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{runs: make(map[string]Run)}
}

func (r *MemoryRepository) Create(ctx context.Context, execution Run) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	snapshot := Clone(execution)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.runs[execution.ID]; exists {
		return fmt.Errorf("%w: %s", ErrRunAlreadyExists, execution.ID)
	}
	r.runs[execution.ID] = snapshot
	return nil
}

func (r *MemoryRepository) Save(ctx context.Context, execution Run) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	snapshot := Clone(execution)
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.runs[execution.ID]; !exists {
		return fmt.Errorf("%w: %s", ErrRunNotFound, execution.ID)
	}
	r.runs[execution.ID] = snapshot
	return nil
}

func (r *MemoryRepository) GetByID(ctx context.Context, id string) (Run, error) {
	if err := ctx.Err(); err != nil {
		return Run{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	execution, exists := r.runs[id]
	if !exists {
		return Run{}, fmt.Errorf("%w: %s", ErrRunNotFound, id)
	}
	return Clone(execution), nil
}
