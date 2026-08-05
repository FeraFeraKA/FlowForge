package run

import (
	"slices"
	"time"

	"flowforge/backend/internal/pipeline"
)

func Clone(source Run) Run {
	clone := source
	clone.JobDefinitions = cloneDefinitions(source.JobDefinitions)
	clone.Jobs = slices.Clone(source.Jobs)
	clone.StartedAt = cloneTime(source.StartedAt)
	clone.FinishedAt = cloneTime(source.FinishedAt)
	for index := range clone.Jobs {
		clone.Jobs[index].StartedAt = cloneTime(source.Jobs[index].StartedAt)
		clone.Jobs[index].FinishedAt = cloneTime(source.Jobs[index].FinishedAt)
	}
	return clone
}

func cloneDefinitions(source []pipeline.JobDefinition) []pipeline.JobDefinition {
	clone := slices.Clone(source)
	for index := range clone {
		clone[index].DependsOn = slices.Clone(source[index].DependsOn)
		if source[index].DelayConfig != nil {
			config := *source[index].DelayConfig
			clone[index].DelayConfig = &config
		}
	}
	return clone
}

func cloneTime(source *time.Time) *time.Time {
	if source == nil {
		return nil
	}
	clone := *source
	return &clone
}
