package httpapi

import "time"

type runResponse struct {
	ID           string           `json:"id"`
	PipelineID   string           `json:"pipelineId"`
	PipelineName string           `json:"pipelineName"`
	Status       string           `json:"status"`
	MaxParallel  int              `json:"maxParallel"`
	Jobs         []jobRunResponse `json:"jobs"`
	CreatedAt    time.Time        `json:"createdAt"`
	StartedAt    *time.Time       `json:"startedAt"`
	FinishedAt   *time.Time       `json:"finishedAt"`
}

type jobRunResponse struct {
	JobID      string               `json:"jobId"`
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	DependsOn  []string             `json:"dependsOn"`
	Timeout    *string              `json:"timeout,omitempty"`
	MaxRetries int                  `json:"maxRetries"`
	RetryDelay string               `json:"retryDelay"`
	Config     *delayConfigResponse `json:"config,omitempty"`
	Status     string               `json:"status"`
	Attempt    int                  `json:"attempt"`
	Progress   int                  `json:"progress"`
	Error      *string              `json:"error"`
	StartedAt  *time.Time           `json:"startedAt"`
	FinishedAt *time.Time           `json:"finishedAt"`
}
