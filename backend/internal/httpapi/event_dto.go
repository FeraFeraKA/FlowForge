package httpapi

import (
	"time"

	"flowforge/backend/internal/events"
	runmodel "flowforge/backend/internal/run"
)

type eventResponse struct {
	ID        string      `json:"id"`
	Type      events.Type `json:"type"`
	RunID     string      `json:"runId"`
	JobID     string      `json:"jobId,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	Payload   any         `json:"payload"`
}

type jobEventPayload struct {
	Job jobEventResponse `json:"job"`
}

type jobEventResponse struct {
	JobID      string     `json:"jobId"`
	Status     string     `json:"status"`
	Attempt    int        `json:"attempt"`
	Progress   int        `json:"progress"`
	Error      *string    `json:"error"`
	StartedAt  *time.Time `json:"startedAt"`
	FinishedAt *time.Time `json:"finishedAt"`
}

func toEventResponse(event events.Event) eventResponse {
	result := eventResponse{
		ID: event.ID, Type: event.Type, RunID: event.RunID,
		JobID: event.JobID, Timestamp: event.Timestamp, Payload: event.Payload,
	}
	switch payload := event.Payload.(type) {
	case events.RunPayload:
		result.Payload = struct {
			Run runResponse `json:"run"`
		}{Run: toRunResponse(payload.Run)}
	case events.JobPayload:
		result.Payload = jobEventPayload{Job: toJobEventResponse(payload.Job)}
	}
	return result
}

func toJobEventResponse(job runmodel.JobRun) jobEventResponse {
	result := jobEventResponse{
		JobID: job.JobID, Status: string(job.Status), Attempt: job.Attempt,
		Progress: job.Progress, StartedAt: job.StartedAt, FinishedAt: job.FinishedAt,
	}
	if job.Error != "" {
		result.Error = &job.Error
	}
	return result
}
