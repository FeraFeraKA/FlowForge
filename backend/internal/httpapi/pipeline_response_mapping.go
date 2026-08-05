package httpapi

import (
	"slices"

	"flowforge/backend/internal/pipeline"
)

func toPipelineResponse(source pipeline.Pipeline) pipelineResponse {
	jobs := make([]jobDefinitionResponse, len(source.Jobs))
	for index, definition := range source.Jobs {
		job := jobDefinitionResponse{
			ID: definition.ID, Name: definition.Name, Type: string(definition.Type),
			DependsOn: slices.Clone(definition.DependsOn), MaxRetries: definition.MaxRetries,
			RetryDelay: definition.RetryDelay.String(),
		}
		if definition.Timeout > 0 {
			timeout := definition.Timeout.String()
			job.Timeout = &timeout
		}
		if definition.DelayConfig != nil {
			job.Config = &delayConfigResponse{Duration: definition.DelayConfig.Duration.String()}
		}
		jobs[index] = job
	}
	return pipelineResponse{
		ID: source.ID, Name: source.Name, MaxParallel: source.MaxParallel,
		Jobs: jobs, CreatedAt: source.CreatedAt,
	}
}

func toPipelineListResponse(sources []pipeline.Pipeline) pipelineListResponse {
	items := make([]pipelineResponse, len(sources))
	for index, source := range sources {
		items[index] = toPipelineResponse(source)
	}
	return pipelineListResponse{Items: items, Total: len(items)}
}
