package scheduler

import (
	"fmt"

	runmodel "flowforge/backend/internal/run"
)

func (s *SequentialScheduler) blockImpossibleJobs(state executionState) (bool, error) {
	blockedAny := false

	for {
		blockedThisPass := false
		for _, definition := range state.execution.JobDefinitions {
			job := state.jobs[definition.ID]
			if job.Status != runmodel.JobStatusPending {
				continue
			}

			shouldBlock := false
			for _, dependencyID := range definition.DependsOn {
				dependency, exists := state.jobs[dependencyID]
				if !exists {
					return false, fmt.Errorf(
						"%w: job %q has unknown dependency %q",
						ErrInvalidRun,
						definition.ID,
						dependencyID,
					)
				}
				if dependency.Status == runmodel.JobStatusFailed ||
					dependency.Status == runmodel.JobStatusBlocked ||
					dependency.Status == runmodel.JobStatusCancelled {
					shouldBlock = true
					break
				}
			}

			if !shouldBlock {
				continue
			}
			if err := job.TransitionTo(runmodel.JobStatusBlocked); err != nil {
				return false, fmt.Errorf("block job %q: %w", job.JobID, err)
			}
			finishedAt := s.now()
			job.FinishedAt = &finishedAt
			blockedThisPass = true
			blockedAny = true
		}

		if !blockedThisPass {
			return blockedAny, nil
		}
	}
}
