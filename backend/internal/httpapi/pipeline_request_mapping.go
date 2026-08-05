package httpapi

import (
	"fmt"
	"slices"
	"time"

	"flowforge/backend/internal/pipeline"
)

func (request createPipelineRequest) params() (pipeline.CreateParams, error) {
	definitions := make([]pipeline.JobDefinition, len(request.Jobs))
	errs := make(pipeline.ValidationErrors, 0)
	for index, job := range request.Jobs {
		prefix := fmt.Sprintf("jobs[%d]", index)
		timeout, validationErr := parseDuration(job.Timeout, prefix+".timeout", false)
		if validationErr != nil {
			errs = append(errs, *validationErr)
		}
		retryDelay, validationErr := parseDuration(job.RetryDelay, prefix+".retryDelay", false)
		if validationErr != nil {
			errs = append(errs, *validationErr)
		}

		definition := pipeline.JobDefinition{
			ID: job.ID, Name: job.Name, Type: pipeline.ExecutorType(job.Type),
			DependsOn: slices.Clone(job.DependsOn), Timeout: timeout,
			MaxRetries: job.MaxRetries, RetryDelay: retryDelay,
		}
		if job.Config != nil {
			duration, validationErr := parseDuration(job.Config.Duration, prefix+".config.duration", true)
			if validationErr != nil {
				errs = append(errs, *validationErr)
			}
			definition.DelayConfig = &pipeline.DelayConfig{Duration: duration}
		}
		definitions[index] = definition
	}
	if len(errs) > 0 {
		return pipeline.CreateParams{}, errs
	}
	return pipeline.CreateParams{
		Name: request.Name, MaxParallel: request.MaxParallel, Jobs: definitions,
	}, nil
}

func parseDuration(value, field string, required bool) (time.Duration, *pipeline.ValidationError) {
	if value == "" && !required {
		return 0, nil
	}
	duration, err := time.ParseDuration(value)
	if err == nil {
		return duration, nil
	}
	return 0, &pipeline.ValidationError{
		Field: field, Code: "invalid_duration",
		Message: "duration must use a value such as 500ms, 2s, or 1m30s",
	}
}
