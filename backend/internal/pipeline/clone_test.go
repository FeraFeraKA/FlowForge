package pipeline

import (
	"testing"
	"time"
)

func TestClonePipelineDoesNotShareMutableData(t *testing.T) {
	original := validPipeline()
	original.Jobs[0].DependsOn = []string{"install"}

	clone := clonePipeline(original)
	original.Jobs[0].Name = "Changed original"
	original.Jobs[0].DependsOn[0] = "changed-original"
	original.Jobs[0].DelayConfig.Duration = 2 * time.Second

	if clone.Jobs[0].Name == original.Jobs[0].Name {
		t.Error("clone shares Jobs backing array with original")
	}
	if clone.Jobs[0].DependsOn[0] == original.Jobs[0].DependsOn[0] {
		t.Error("clone shares DependsOn backing array with original")
	}
	if clone.Jobs[0].DelayConfig.Duration == original.Jobs[0].DelayConfig.Duration {
		t.Error("clone shares DelayConfig pointer with original")
	}

	clone.Jobs[0].Name = "Changed clone"
	if original.Jobs[0].Name == clone.Jobs[0].Name {
		t.Error("original changed after mutating clone")
	}
}
