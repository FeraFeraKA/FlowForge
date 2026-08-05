package httpapi

import (
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"flowforge/backend/internal/events"
)

func TestWriteSSE(t *testing.T) {
	response := httptest.NewRecorder()
	event := events.Event{
		ID: "7", Type: events.JobCompleted, RunID: "run-1", JobID: "tests",
		Timestamp: time.Date(2026, time.August, 5, 10, 0, 0, 0, time.UTC),
	}

	if err := writeSSE(response, event); err != nil {
		t.Fatalf("write SSE: %v", err)
	}

	body := response.Body.String()
	for _, fragment := range []string{
		"id: 7\n", "event: job.completed\n", `"runId":"run-1"`, "\n\n",
	} {
		if !strings.Contains(body, fragment) {
			t.Errorf("SSE frame does not contain %q: %q", fragment, body)
		}
	}
}
