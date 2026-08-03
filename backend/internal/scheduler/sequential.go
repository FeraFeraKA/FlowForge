package scheduler

import (
	"errors"
	"fmt"
	"time"

	"flowforge/backend/internal/pipeline"
)

var (
	ErrInvalidRun            = errors.New("invalid run snapshot")
	ErrExecutorNotRegistered = errors.New("executor not registered")
	ErrSchedulerStalled      = errors.New("scheduler stalled")
)

type SequentialScheduler struct {
	executors map[pipeline.ExecutorType]Executor
	now       func() time.Time
}

func NewSequentialScheduler(executors map[pipeline.ExecutorType]Executor) *SequentialScheduler {
	registry := make(map[pipeline.ExecutorType]Executor, len(executors))
	for executorType, executor := range executors {
		registry[executorType] = executor
	}

	return &SequentialScheduler{
		executors: registry,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

func (s *SequentialScheduler) executorFor(executorType pipeline.ExecutorType) (Executor, error) {
	executor, exists := s.executors[executorType]
	if !exists || executor == nil {
		return nil, fmt.Errorf("%w: %s", ErrExecutorNotRegistered, executorType)
	}
	return executor, nil
}
