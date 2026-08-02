package pipeline

import (
	"testing"
	"time"
)

func TestDependencyCycleValidation(t *testing.T) {
	tests := []struct {
		name      string
		graph     map[string][]string
		wantCycle bool
	}{
		{
			name: "diamond DAG",
			graph: map[string][]string{
				"install": nil,
				"tests":   {"install"},
				"build":   {"install"},
				"deploy":  {"tests", "build"},
			},
		},
		{
			name:      "two-job cycle",
			graph:     map[string][]string{"first": {"second"}, "second": {"first"}},
			wantCycle: true,
		},
		{
			name: "long cycle",
			graph: map[string][]string{
				"first": {"second"}, "second": {"third"}, "third": {"first"},
			},
			wantCycle: true,
		},
		{
			name: "cycle in disconnected component",
			graph: map[string][]string{
				"root": nil, "child": {"root"}, "left": {"right"}, "right": {"left"},
			},
			wantCycle: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			pipeline := pipelineFromGraph(test.graph)
			err := pipeline.Validate()
			if !test.wantCycle {
				if err != nil {
					t.Fatalf("Validate() error = %v, want nil", err)
				}
				return
			}

			errs := requireValidationErrors(t, err)
			if !containsValidationError(errs, "jobs", "dependency_cycle") {
				t.Errorf("validation errors = %#v, want dependency_cycle", errs)
			}
		})
	}
}

func pipelineFromGraph(graph map[string][]string) Pipeline {
	pipeline := Pipeline{Name: "Graph", MaxParallel: 2, Jobs: make([]JobDefinition, 0, len(graph))}
	for jobID, dependencies := range graph {
		pipeline.Jobs = append(pipeline.Jobs, JobDefinition{
			ID: jobID, Name: jobID, Type: ExecutorTypeDelay, DependsOn: dependencies,
			DelayConfig: &DelayConfig{Duration: time.Second},
		})
	}
	return pipeline
}
