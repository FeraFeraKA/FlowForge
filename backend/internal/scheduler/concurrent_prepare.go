package scheduler

import (
	"context"
	"fmt"

	runmodel "flowforge/backend/internal/run"
)

func (s *ConcurrentScheduler) prepare(
	ctx context.Context,
	execution *runmodel.Run,
) (*concurrentRuntime, error) {
	state, err := indexExecution(execution)
	if err != nil {
		return nil, err
	}
	if execution.Status != runmodel.RunStatusPending {
		return nil, fmt.Errorf("%w: run starts as %q", ErrInvalidRun, execution.Status)
	}
	if execution.MaxParallel < 1 {
		return nil, fmt.Errorf("%w: maxParallel must be positive", ErrInvalidRun)
	}
	for _, job := range execution.Jobs {
		if job.Status != runmodel.JobStatusPending {
			return nil, fmt.Errorf("%w: job %q starts as %q", ErrInvalidRun, job.JobID, job.Status)
		}
	}
	for _, definition := range execution.JobDefinitions {
		if definition.MaxRetries < 0 || definition.RetryDelay < 0 || definition.Timeout < 0 {
			return nil, fmt.Errorf("%w: job %q has negative execution settings", ErrInvalidRun, definition.ID)
		}
		if _, err := s.base.executorFor(definition.Type); err != nil {
			return nil, err
		}
	}

	workCtx, cancel := context.WithCancel(ctx)
	return &concurrentRuntime{
		ctx:     workCtx,
		cancel:  cancel,
		state:   state,
		results: make(chan attemptMessage, execution.MaxParallel),
		retries: make(chan retryMessage, len(execution.Jobs)),
	}, nil
}
