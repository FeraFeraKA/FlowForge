package scheduler

import (
	"fmt"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

type executionState struct {
	execution   *runmodel.Run
	definitions map[string]pipeline.JobDefinition
	jobs        map[string]*runmodel.JobRun
}

func indexExecution(execution *runmodel.Run) (executionState, error) {
	if execution == nil {
		return executionState{}, fmt.Errorf("%w: run is nil", ErrInvalidRun)
	}
	if len(execution.JobDefinitions) == 0 {
		return executionState{}, fmt.Errorf("%w: run has no job definitions", ErrInvalidRun)
	}
	if len(execution.JobDefinitions) != len(execution.Jobs) {
		return executionState{}, fmt.Errorf(
			"%w: %d definitions but %d job states",
			ErrInvalidRun,
			len(execution.JobDefinitions),
			len(execution.Jobs),
		)
	}

	state := executionState{
		execution:   execution,
		definitions: make(map[string]pipeline.JobDefinition, len(execution.JobDefinitions)),
		jobs:        make(map[string]*runmodel.JobRun, len(execution.Jobs)),
	}
	for _, definition := range execution.JobDefinitions {
		if _, duplicate := state.definitions[definition.ID]; duplicate {
			return executionState{}, fmt.Errorf("%w: duplicate definition %q", ErrInvalidRun, definition.ID)
		}
		state.definitions[definition.ID] = definition
	}

	for index := range execution.Jobs {
		job := &execution.Jobs[index]
		if _, duplicate := state.jobs[job.JobID]; duplicate {
			return executionState{}, fmt.Errorf("%w: duplicate job state %q", ErrInvalidRun, job.JobID)
		}
		if _, exists := state.definitions[job.JobID]; !exists {
			return executionState{}, fmt.Errorf("%w: state for unknown job %q", ErrInvalidRun, job.JobID)
		}
		state.jobs[job.JobID] = job
	}

	for jobID := range state.definitions {
		if _, exists := state.jobs[jobID]; !exists {
			return executionState{}, fmt.Errorf("%w: state for job %q is missing", ErrInvalidRun, jobID)
		}
	}

	return state, nil
}
