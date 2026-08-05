package scheduler

import (
	"context"
	"fmt"

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

func (s *ConcurrentScheduler) Execute(ctx context.Context, execution *runmodel.Run) error {
	runtime, err := s.prepare(ctx, execution)
	if err != nil {
		return err
	}
	if ctxErr := runtime.ctx.Err(); ctxErr != nil {
		return s.stop(runtime, ctxErr)
	}
	if err := execution.TransitionTo(runmodel.RunStatusRunning); err != nil {
		return s.stop(runtime, fmt.Errorf("start run %q: %w", execution.ID, err))
	}
	startedAt := s.base.now()
	execution.StartedAt = &startedAt
	if err := s.notify(execution); err != nil {
		return s.stop(runtime, fmt.Errorf("observe started run: %w", err))
	}

	for {
		if ctxErr := runtime.ctx.Err(); ctxErr != nil {
			return s.stop(runtime, ctxErr)
		}
		blocked, err := s.base.blockImpossibleJobs(runtime.state)
		if err != nil {
			return s.stop(runtime, err)
		}
		if blocked {
			if err := s.notify(execution); err != nil {
				return s.stop(runtime, fmt.Errorf("observe blocked jobs: %w", err))
			}
		}
		if _, err := s.scheduleAvailable(runtime); err != nil {
			return s.stop(runtime, err)
		}
		if runtime.active == 0 && runtime.retryWaiting == 0 {
			return s.complete(runtime)
		}

		select {
		case <-runtime.ctx.Done():
			return s.stop(runtime, runtime.ctx.Err())
		case message := <-runtime.results:
			if err := s.handleAttemptResult(runtime, message); err != nil {
				return s.stop(runtime, err)
			}
		case message := <-runtime.retries:
			if err := s.handleRetry(runtime, message); err != nil {
				return s.stop(runtime, err)
			}
		}
	}
}
