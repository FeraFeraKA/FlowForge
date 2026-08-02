package pipeline

import "testing"

func TestDependencyValidation(t *testing.T) {
	tests := []struct {
		name   string
		change func(*Pipeline)
		field  string
		code   string
	}{
		{
			name: "duplicate job ID",
			change: func(p *Pipeline) {
				duplicate := p.Jobs[0]
				p.Jobs = append(p.Jobs, duplicate)
			},
			field: "jobs[1].id", code: "duplicate",
		},
		{
			name: "unknown dependency",
			change: func(p *Pipeline) {
				p.Jobs[0].DependsOn = []string{"missing"}
			},
			field: "jobs[0].dependsOn[0]", code: "unknown_dependency",
		},
		{
			name: "self dependency",
			change: func(p *Pipeline) {
				p.Jobs[0].DependsOn = []string{"build"}
			},
			field: "jobs[0].dependsOn[0]", code: "self_dependency",
		},
		{
			name: "duplicate dependency",
			change: func(p *Pipeline) {
				dependency := p.Jobs[0]
				dependency.ID, dependency.Name = "install", "Install"
				p.Jobs = append(p.Jobs, dependency)
				p.Jobs[0].DependsOn = []string{"install", "install"}
			},
			field: "jobs[0].dependsOn[1]", code: "duplicate",
		},
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
