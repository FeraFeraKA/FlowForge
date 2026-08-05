package scheduler

import (
	"fmt"
	"time"

	runmodel "flowforge/backend/internal/run"
)

func (s *ConcurrentScheduler) handleAttemptResult(
	runtime *concurrentRuntime,
	message attemptMessage,
) error {
	runtime.active--
	job, exists := runtime.state.jobs[message.jobID]
	if !exists {
		return fmt.Errorf("%w: result for unknown job %q", ErrInvalidRun, message.jobID)
	}
	if runtime.ctx.Err() != nil {
		return nil
	}

	if message.err == nil {
		if err := job.TransitionTo(runmodel.JobStatusCompleted); err != nil {
			return fmt.Errorf("complete job %q: %w", job.JobID, err)
		}
		job.Progress = 100
		s.base.finishJob(job)
		return s.notify(runtime.state.execution)
	}

	if err := job.TransitionTo(runmodel.JobStatusFailed); err != nil {
		return fmt.Errorf("fail job %q: %w", job.JobID, err)
	}
	job.Error = message.err.Error()
	definition := runtime.state.definitions[job.JobID]
	if job.Attempt > definition.MaxRetries {
		s.base.finishJob(job)
		return s.notify(runtime.state.execution)
	}

	if err := job.TransitionTo(runmodel.JobStatusRetrying); err != nil {
		return fmt.Errorf("retry job %q: %w", job.JobID, err)
	}
	if err := s.notify(runtime.state.execution); err != nil {
		return err
	}
	s.launchRetryWait(runtime, job.JobID, definition.RetryDelay)
	return nil
}

func (s *ConcurrentScheduler) launchRetryWait(
	runtime *concurrentRuntime,
	jobID string,
	delay time.Duration,
) {
	runtime.retryWaiting++
	runtime.workers.Add(1)
	go func() {
		defer runtime.workers.Done()
		err := waitForRetry(runtime.ctx, delay)
		runtime.retries <- retryMessage{jobID: jobID, err: err}
	}()
}
