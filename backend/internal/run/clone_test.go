package run

import (
	"testing"
	"time"

	"flowforge/backend/internal/pipeline"
)

func TestCloneDoesNotShareMutableData(t *testing.T) {
	startedAt := time.Now()
	source := Run{
		JobDefinitions: []pipeline.JobDefinition{{
			ID: "job", DependsOn: []string{"dependency"},
			DelayConfig: &pipeline.DelayConfig{Duration: time.Second},
		}},
		Jobs:      []JobRun{{JobID: "job", StartedAt: &startedAt}},
		StartedAt: &startedAt,
	}

	clone := Clone(source)
	source.JobDefinitions[0].DependsOn[0] = "changed"
	source.JobDefinitions[0].DelayConfig.Duration = 2 * time.Second
	source.Jobs[0].JobID = "changed"
	*source.StartedAt = source.StartedAt.Add(time.Hour)

	if clone.JobDefinitions[0].DependsOn[0] == "changed" ||
		clone.JobDefinitions[0].DelayConfig.Duration == 2*time.Second ||
		clone.Jobs[0].JobID == "changed" || clone.StartedAt.Equal(*source.StartedAt) {
		t.Error("Clone() shares mutable data with source")
	}
}
