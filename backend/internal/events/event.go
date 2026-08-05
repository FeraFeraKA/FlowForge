package events

import "time"

type Type string

const (
	RunStarted   Type = "run.started"
	RunCompleted Type = "run.completed"
	RunFailed    Type = "run.failed"
	RunCancelled Type = "run.cancelled"
	JobQueued    Type = "job.queued"
	JobStarted   Type = "job.started"
	JobProgress  Type = "job.progress"
	JobRetrying  Type = "job.retrying"
	JobCompleted Type = "job.completed"
	JobFailed    Type = "job.failed"
	JobBlocked   Type = "job.blocked"
	JobCancelled Type = "job.cancelled"
)

type Event struct {
	ID        string    `json:"id"`
	Type      Type      `json:"type"`
	RunID     string    `json:"runId"`
	JobID     string    `json:"jobId,omitempty"`
	Timestamp time.Time `json:"timestamp"`
	Payload   any       `json:"payload"`
}
