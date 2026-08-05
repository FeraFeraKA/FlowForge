package run

import (
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
)

func TestNewCreatesIndependentPendingRun(t *testing.T) {
	createdAt := time.Now()
	source := pipeline.Pipeline{
		ID: "pipeline", Name: "Build", MaxParallel: 2,
		Jobs: []pipeline.JobDefinition{{ID: "job", DependsOn: []string{"dependency"}}},
	}

	execution := New("run", source, createdAt)
	source.Jobs[0].DependsOn[0] = "changed"

	if execution.ID != "run" || execution.PipelineID != "pipeline" ||
		execution.Status != RunStatusPending || execution.MaxParallel != 2 {
		t.Errorf("New() returned unexpected run: %#v", execution)
	}
	if len(execution.Jobs) != 1 || execution.Jobs[0].Status != JobStatusPending {
		t.Errorf("New() jobs = %#v, want one pending job", execution.Jobs)
	}
	if execution.JobDefinitions[0].DependsOn[0] == "changed" {
		t.Error("New() shares definitions with pipeline")
	}
}
