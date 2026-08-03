package scheduler

import (
	"fmt"

	runmodel "flowforge/backend/internal/run"
)

func (s *SequentialScheduler) cancelExecution(state executionState, cause error) error {
	for index := range state.execution.Jobs {
		job := &state.execution.Jobs[index]
		switch job.Status {
		case runmodel.JobStatusPending,
			runmodel.JobStatusQueued,
			runmodel.JobStatusRunning,
			runmodel.JobStatusRetrying:
			if err := job.TransitionTo(runmodel.JobStatusCancelled); err != nil {
				return fmt.Errorf("cancel job %q: %w", job.JobID, err)
			}
			if cause != nil {
				job.Error = cause.Error()
			}
			s.finishJob(job)
		}
	}

	if state.execution.Status != runmodel.RunStatusCancelled {
		if err := state.execution.TransitionTo(runmodel.RunStatusCancelled); err != nil {
			return fmt.Errorf("cancel run %q: %w", state.execution.ID, err)
		}
	}
	finishedAt := s.now()
	state.execution.FinishedAt = &finishedAt
	return nil
}
