package scheduler

import (
	"context"
	"fmt"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

type attemptResult struct {
	executionError error
}

func (s *SequentialScheduler) executeAttempt(
	ctx context.Context,
	job *runmodel.JobRun,
	definition pipeline.JobDefinition,
) (attemptResult, error) {
	executor, err := s.executorFor(definition.Type)
	if err != nil {
		return attemptResult{}, err
	}
	if definition.Timeout < 0 {
		return attemptResult{}, fmt.Errorf("%w: job %q has negative timeout", ErrInvalidRun, definition.ID)
	}
	if err := job.TransitionTo(runmodel.JobStatusRunning); err != nil {
		return attemptResult{}, fmt.Errorf("start job %q: %w", job.JobID, err)
	}

	job.Attempt++
	job.Progress = 0
	job.Error = ""
	if job.StartedAt == nil {
		startedAt := s.now()
		job.StartedAt = &startedAt
	}

	attemptCtx := ctx
	var cancel context.CancelFunc
	if definition.Timeout > 0 {
		attemptCtx, cancel = context.WithTimeout(ctx, definition.Timeout)
	}

	executionErr := executor.Execute(attemptCtx, definition)
	if cancel != nil {
		cancel()
	}

	return attemptResult{executionError: executionErr}, nil
}
