package scheduler

import (
	"fmt"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func (state executionState) nextReady() (*runmodel.JobRun, pipeline.JobDefinition, bool, error) {
	for _, definition := range state.execution.JobDefinitions {
		job := state.jobs[definition.ID]
		if job.Status != runmodel.JobStatusPending {
			continue
		}

		ready := true
		for _, dependencyID := range definition.DependsOn {
			dependency, exists := state.jobs[dependencyID]
			if !exists {
				return nil, pipeline.JobDefinition{}, false, fmt.Errorf(
					"%w: job %q has unknown dependency %q",
					ErrInvalidRun,
					definition.ID,
					dependencyID,
				)
			}
			if dependency.Status != runmodel.JobStatusCompleted {
				ready = false
				break
			}
		}

		if ready {
			return job, definition, true, nil
		}
	}

	return nil, pipeline.JobDefinition{}, false, nil
}
