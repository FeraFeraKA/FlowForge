package scheduler

import (
	"fmt"

	runmodel "flowforge/backend/internal/run"
)

func (s *ConcurrentScheduler) handleRetry(
	runtime *concurrentRuntime,
	message retryMessage,
) error {
	runtime.retryWaiting--
	job, exists := runtime.state.jobs[message.jobID]
	if !exists {
		return fmt.Errorf("%w: retry for unknown job %q", ErrInvalidRun, message.jobID)
	}
	if message.err != nil {
		if runtime.ctx.Err() != nil {
			return nil
		}
		return fmt.Errorf("wait to retry job %q: %w", message.jobID, message.err)
	}
	if err := job.TransitionTo(runmodel.JobStatusQueued); err != nil {
		return fmt.Errorf("requeue job %q: %w", job.JobID, err)
	}
	if err := s.notify(runtime.state.execution); err != nil {
		return fmt.Errorf("observe requeued job %q: %w", job.JobID, err)
	}
	return nil
}
