package httpapi

import (
	"slices"

	runmodel "flowforge/backend/internal/run"
)

func toRunResponse(source runmodel.Run) runResponse {
	jobs := make([]jobRunResponse, len(source.Jobs))
	definitions := make(map[string]int, len(source.JobDefinitions))
	for index, definition := range source.JobDefinitions {
		definitions[definition.ID] = index
	}
	for index, state := range source.Jobs {
		definition := source.JobDefinitions[definitions[state.JobID]]
		job := jobRunResponse{
			JobID: state.JobID, Name: definition.Name, Type: string(definition.Type),
			DependsOn: slices.Clone(definition.DependsOn), MaxRetries: definition.MaxRetries,
			RetryDelay: definition.RetryDelay.String(), Status: string(state.Status),
			Attempt: state.Attempt, Progress: state.Progress,
			StartedAt: state.StartedAt, FinishedAt: state.FinishedAt,
		}
		if definition.Timeout > 0 {
			timeout := definition.Timeout.String()
			job.Timeout = &timeout
		}
		if definition.DelayConfig != nil {
			job.Config = &delayConfigResponse{Duration: definition.DelayConfig.Duration.String()}
		}
		if state.Error != "" {
			errorMessage := state.Error
			job.Error = &errorMessage
		}
		jobs[index] = job
	}
	return runResponse{
		ID: source.ID, PipelineID: source.PipelineID, PipelineName: source.PipelineName,
		Status: string(source.Status), MaxParallel: source.MaxParallel, Jobs: jobs,
		CreatedAt: source.CreatedAt, StartedAt: source.StartedAt, FinishedAt: source.FinishedAt,
	}
}
