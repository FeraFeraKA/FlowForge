package scheduler

import (
	"fmt"

	runmodel "flowforge/backend/internal/run"
)

func (s *SequentialScheduler) finishExecution(state executionState) error {
	hasFailure := false
	for _, job := range state.execution.Jobs {
		switch job.Status {
		case runmodel.JobStatusCompleted:
		case runmodel.JobStatusFailed, runmodel.JobStatusBlocked:
			hasFailure = true
		default:
			return fmt.Errorf(
				"%w: job %q remains %q",
				ErrSchedulerStalled,
				job.JobID,
				job.Status,
			)
		}
	}

	nextStatus := runmodel.RunStatusCompleted
	if hasFailure {
		nextStatus = runmodel.RunStatusFailed
	}
	if err := state.execution.TransitionTo(nextStatus); err != nil {
		return fmt.Errorf("finish run %q: %w", state.execution.ID, err)
	}
	finishedAt := s.now()
	state.execution.FinishedAt = &finishedAt
	return nil
}
