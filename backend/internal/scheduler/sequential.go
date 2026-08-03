package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
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

func (s *SequentialScheduler) Execute(ctx context.Context, execution *runmodel.Run) error {
	state, err := indexExecution(execution)
	if err != nil {
		return err
	}
	if execution.Status != runmodel.RunStatusPending {
		return fmt.Errorf("%w: run starts as %q", ErrInvalidRun, execution.Status)
	}
	for _, job := range execution.Jobs {
		if job.Status != runmodel.JobStatusPending {
			return fmt.Errorf("%w: job %q starts as %q", ErrInvalidRun, job.JobID, job.Status)
		}
	}
	for _, definition := range execution.JobDefinitions {
		if _, err := s.executorFor(definition.Type); err != nil {
			return err
		}
	}

	if err := ctx.Err(); err != nil {
		if cancelErr := s.cancelExecution(state, err); cancelErr != nil {
			return cancelErr
		}
		return err
	}
	if err := execution.TransitionTo(runmodel.RunStatusRunning); err != nil {
		return fmt.Errorf("start run %q: %w", execution.ID, err)
	}
	startedAt := s.now()
	execution.StartedAt = &startedAt

	for {
		if ctxErr := ctx.Err(); ctxErr != nil {
			if err := s.cancelExecution(state, ctxErr); err != nil {
				return err
			}
			return ctxErr
		}
		if _, err := s.blockImpossibleJobs(state); err != nil {
			return err
		}

		job, definition, ready, err := state.nextReady()
		if err != nil {
			return err
		}
		if !ready {
			return s.finishExecution(state)
		}

		if err := s.runJob(ctx, job, definition); err != nil {
			if ctxErr := ctx.Err(); ctxErr != nil && errors.Is(err, ctxErr) {
				if cancelErr := s.cancelExecution(state, ctxErr); cancelErr != nil {
					return cancelErr
				}
				return ctxErr
			}
			return fmt.Errorf("run job %q: %w", job.JobID, err)
		}
	}
}
