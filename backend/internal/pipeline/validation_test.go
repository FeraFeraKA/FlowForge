package pipeline

import (
	"testing"
	"time"
)

func TestPipelineFieldValidation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Pipeline)
		field  string
		code   string
	}{
		{"missing pipeline name", func(p *Pipeline) { p.Name = " " }, "name", "required"},
		{"invalid max parallel", func(p *Pipeline) { p.MaxParallel = 0 }, "maxParallel", "must_be_positive"},
		{"missing jobs", func(p *Pipeline) { p.Jobs = nil }, "jobs", "required"},
		{"invalid job ID", func(p *Pipeline) { p.Jobs[0].ID = "bad id" }, "jobs[0].id", "invalid_format"},
		{"missing job name", func(p *Pipeline) { p.Jobs[0].Name = "" }, "jobs[0].name", "required"},
		{"negative retries", func(p *Pipeline) { p.Jobs[0].MaxRetries = -1 }, "jobs[0].maxRetries", "must_not_be_negative"},
		{"missing delay config", func(p *Pipeline) { p.Jobs[0].DelayConfig = nil }, "jobs[0].config", "required"},
		{"invalid duration", func(p *Pipeline) { p.Jobs[0].DelayConfig.Duration = 0 }, "jobs[0].config.duration", "must_be_positive"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pipeline := validPipeline()
			test.change(&pipeline)

			errs := requireValidationErrors(t, pipeline.Validate())
			if !containsValidationError(errs, test.field, test.code) {
				t.Errorf("validation errors = %#v, want field %q with code %q", errs, test.field, test.code)
			}
		})
	}
}

func validPipeline() Pipeline {
	return Pipeline{
		ID:          "pipeline-1",
		Name:        "Build",
		MaxParallel: 1,
		Jobs: []JobDefinition{{
			ID:          "build",
			Name:        "Build",
			Type:        ExecutorTypeDelay,
			DelayConfig: &DelayConfig{Duration: time.Second},
		}},
	}
}

func requireValidationErrors(t *testing.T, err error) ValidationErrors {
	t.Helper()
	if err == nil {
		t.Fatal("Validate() error = nil, want validation error")
	}
	errs, ok := err.(ValidationErrors)
	if !ok {
		t.Fatalf("Validate() error type = %T, want ValidationErrors", err)
	}
	return errs
}

func containsValidationError(errs ValidationErrors, field, code string) bool {
	for _, err := range errs {
		if err.Field == field && err.Code == code {
			return true
		}
	}
	return false
}
