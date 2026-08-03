package scheduler

import (
	"context"
	"fmt"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func (s *SequentialScheduler) runJob(
	ctx context.Context,
	job *runmodel.JobRun,
	definition pipeline.JobDefinition,
) error {
	if definition.MaxRetries < 0 || definition.RetryDelay < 0 || definition.Timeout < 0 {
		return fmt.Errorf("%w: job %q has negative execution settings", ErrInvalidRun, definition.ID)
	}
	if err := job.TransitionTo(runmodel.JobStatusQueued); err != nil {
		return fmt.Errorf("queue job %q: %w", job.JobID, err)
	}

	for {
		result, err := s.executeAttempt(ctx, job, definition)
		if err != nil {
			return err
		}
		if result.executionError == nil {
			if err := job.TransitionTo(runmodel.JobStatusCompleted); err != nil {
				return fmt.Errorf("complete job %q: %w", job.JobID, err)
			}
			job.Progress = 100
			s.finishJob(job)
			return nil
		}

		if ctx.Err() != nil {
			if err := job.TransitionTo(runmodel.JobStatusCancelled); err != nil {
				return fmt.Errorf("cancel running job %q: %w", job.JobID, err)
			}
			job.Error = ctx.Err().Error()
			s.finishJob(job)
			return ctx.Err()
		}

		if err := job.TransitionTo(runmodel.JobStatusFailed); err != nil {
			return fmt.Errorf("fail job %q: %w", job.JobID, err)
		}
		job.Error = result.executionError.Error()
		if job.Attempt > definition.MaxRetries {
			s.finishJob(job)
			return nil
		}

		if err := job.TransitionTo(runmodel.JobStatusRetrying); err != nil {
			return fmt.Errorf("retry job %q: %w", job.JobID, err)
		}
		if err := waitForRetry(ctx, definition.RetryDelay); err != nil {
			if ctx.Err() == nil {
				return fmt.Errorf("wait to retry job %q: %w", job.JobID, err)
			}
			if transitionErr := job.TransitionTo(runmodel.JobStatusCancelled); transitionErr != nil {
				return fmt.Errorf("cancel retrying job %q: %w", job.JobID, transitionErr)
			}
			job.Error = ctx.Err().Error()
			s.finishJob(job)
			return ctx.Err()
		}
		if err := job.TransitionTo(runmodel.JobStatusQueued); err != nil {
			return fmt.Errorf("requeue job %q: %w", job.JobID, err)
		}
	}
}

func (s *SequentialScheduler) finishJob(job *runmodel.JobRun) {
	finishedAt := s.now()
	job.FinishedAt = &finishedAt
}
