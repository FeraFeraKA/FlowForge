package scheduler

import (
	"fmt"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func (s *ConcurrentScheduler) scheduleAvailable(runtime *concurrentRuntime) (bool, error) {
	launchedAny := false
	for runtime.active < runtime.state.execution.MaxParallel {
		job, definition, found := queuedJob(runtime.state)
		if !found {
			var err error
			job, definition, found, err = runtime.state.nextReady()
			if err != nil {
				return false, err
			}
			if !found {
				return launchedAny, nil
			}
			if err := job.TransitionTo(runmodel.JobStatusQueued); err != nil {
				return false, fmt.Errorf("queue job %q: %w", job.JobID, err)
			}
			if err := s.notify(runtime.state.execution); err != nil {
				return false, fmt.Errorf("observe queued job %q: %w", job.JobID, err)
			}
		}

		if err := s.launchAttempt(runtime, job, definition); err != nil {
			return false, err
		}
		launchedAny = true
	}
	return launchedAny, nil
}

func queuedJob(state executionState) (*runmodel.JobRun, pipeline.JobDefinition, bool) {
	for _, definition := range state.execution.JobDefinitions {
		job := state.jobs[definition.ID]
		if job.Status == runmodel.JobStatusQueued {
			return job, definition, true
		}
	}
	return nil, pipeline.JobDefinition{}, false
}
