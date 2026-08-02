package pipeline

import "slices"

func clonePipeline(source Pipeline) Pipeline {
	clone := source
	clone.Jobs = slices.Clone(source.Jobs)

	for index := range clone.Jobs {
		clone.Jobs[index].DependsOn = slices.Clone(source.Jobs[index].DependsOn)
		if source.Jobs[index].DelayConfig != nil {
			config := *source.Jobs[index].DelayConfig
			clone.Jobs[index].DelayConfig = &config
		}
	}

	return clone
}
