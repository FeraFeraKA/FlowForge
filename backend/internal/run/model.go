package run

import (
	"time"

	"flowforge/backend/internal/pipeline"
)

type Run struct {
	ID             string
	PipelineID     string
	PipelineName   string
	Status         RunStatus
	MaxParallel    int
	JobDefinitions []pipeline.JobDefinition
	Jobs           []JobRun
	CreatedAt      time.Time
	StartedAt      *time.Time
	FinishedAt     *time.Time
}

type JobRun struct {
	JobID      string
	Status     JobStatus
	Attempt    int
	Progress   int
	Error      string
	StartedAt  *time.Time
	FinishedAt *time.Time
}
