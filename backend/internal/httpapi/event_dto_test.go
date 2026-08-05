package httpapi

import (
	"encoding/json"
	"strings"
	"testing"

	"flowforge/backend/internal/events"
	runmodel "flowforge/backend/internal/run"
)

func TestJobEventUsesPublicJSONContract(t *testing.T) {
	encoded, err := json.Marshal(toEventResponse(events.Event{
		ID: "3", Type: events.JobStarted, RunID: "run-1", JobID: "tests",
		Payload: events.JobPayload{Job: runmodel.JobRun{
			JobID: "tests", Status: runmodel.JobStatusRunning, Attempt: 1,
		}},
	}))
	if err != nil {
		t.Fatalf("marshal event: %v", err)
	}
	text := string(encoded)
	for _, fragment := range []string{
		`"runId":"run-1"`, `"jobId":"tests"`, `"status":"running"`,
	} {
		if !strings.Contains(text, fragment) {
			t.Errorf("event does not contain %s: %s", fragment, text)
		}
	}
	if strings.Contains(text, `"JobID"`) {
		t.Fatalf("event leaked Go field names: %s", text)
	}
}
