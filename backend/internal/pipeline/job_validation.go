package pipeline

import (
	"fmt"
	"strings"
)

func validateJobFields(index int, job JobDefinition) ValidationErrors {
	errs := make(ValidationErrors, 0)
	prefix := fmt.Sprintf("jobs[%d]", index)
	add := func(field, code, message string) {
		errs = append(errs, ValidationError{Field: prefix + field, Code: code, Message: message})
	}

	if job.ID == "" {
		add(".id", "required", "job ID is required")
	} else if !validJobID(job.ID) {
		add(".id", "invalid_format", "job ID may contain only ASCII letters, digits, _ and -")
	}
	if strings.TrimSpace(job.Name) == "" {
		add(".name", "required", "job name is required")
	}
	if job.Timeout < 0 {
		add(".timeout", "must_not_be_negative", "timeout must not be negative")
	}
	if job.MaxRetries < 0 {
		add(".maxRetries", "must_not_be_negative", "maxRetries must not be negative")
	}
	if job.RetryDelay < 0 {
		add(".retryDelay", "must_not_be_negative", "retryDelay must not be negative")
	}

	switch job.Type {
	case ExecutorTypeDelay:
		if job.DelayConfig == nil {
			add(".config", "required", "delay config is required")
		} else if job.DelayConfig.Duration <= 0 {
			add(".config.duration", "must_be_positive", "delay duration must be positive")
		}
	default:
		add(".type", "unsupported", fmt.Sprintf("executor type %q is not supported", job.Type))
	}

	return errs
}

func validJobID(id string) bool {
	for index := 0; index < len(id); index++ {
		character := id[index]
		if character >= 'a' && character <= 'z' ||
			character >= 'A' && character <= 'Z' ||
			character >= '0' && character <= '9' ||
			character == '_' || character == '-' {
			continue
		}
		return false
	}
	return true
}
