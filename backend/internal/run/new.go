package run

import (
	"time"

	"flowforge/backend/internal/pipeline"
)

func New(id string, source pipeline.Pipeline, createdAt time.Time) Run {
	definitions := cloneDefinitions(source.Jobs)
	jobs := make([]JobRun, len(definitions))
	for index, definition := range definitions {
		jobs[index] = JobRun{
			JobID:  definition.ID,
			Status: JobStatusPending,
		}
	}

	return Run{
		ID:             id,
		PipelineID:     source.ID,
		PipelineName:   source.Name,
		Status:         RunStatusPending,
		MaxParallel:    source.MaxParallel,
		JobDefinitions: definitions,
		Jobs:           jobs,
		CreatedAt:      createdAt,
	}
}
