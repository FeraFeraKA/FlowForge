package scheduler

import (
	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

type StateObserver func(snapshot runmodel.Run) error

type ConcurrentScheduler struct {
	base     *SequentialScheduler
	observer StateObserver
}

func NewConcurrentScheduler(
	executors map[pipeline.ExecutorType]Executor,
	observer StateObserver,
) *ConcurrentScheduler {
	return &ConcurrentScheduler{
		base:     NewSequentialScheduler(executors),
		observer: observer,
	}
}

func (s *ConcurrentScheduler) notify(execution *runmodel.Run) error {
	if s.observer == nil {
		return nil
	}
	return s.observer(runmodel.Clone(*execution))
}
