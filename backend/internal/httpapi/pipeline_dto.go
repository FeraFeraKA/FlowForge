package httpapi

import "time"

type createPipelineRequest struct {
	Name        string                 `json:"name"`
	MaxParallel int                    `json:"maxParallel"`
	Jobs        []jobDefinitionRequest `json:"jobs"`
}

type jobDefinitionRequest struct {
	ID         string              `json:"id"`
	Name       string              `json:"name"`
	Type       string              `json:"type"`
	DependsOn  []string            `json:"dependsOn"`
	Timeout    string              `json:"timeout"`
	MaxRetries int                 `json:"maxRetries"`
	RetryDelay string              `json:"retryDelay"`
	Config     *delayConfigRequest `json:"config"`
}

type delayConfigRequest struct {
	Duration string `json:"duration"`
}

type pipelineResponse struct {
	ID          string                  `json:"id"`
	Name        string                  `json:"name"`
	MaxParallel int                     `json:"maxParallel"`
	Jobs        []jobDefinitionResponse `json:"jobs"`
	CreatedAt   time.Time               `json:"createdAt"`
}

type jobDefinitionResponse struct {
	ID         string               `json:"id"`
	Name       string               `json:"name"`
	Type       string               `json:"type"`
	DependsOn  []string             `json:"dependsOn"`
	Timeout    *string              `json:"timeout,omitempty"`
	MaxRetries int                  `json:"maxRetries"`
	RetryDelay string               `json:"retryDelay"`
	Config     *delayConfigResponse `json:"config,omitempty"`
}

type delayConfigResponse struct {
	Duration string `json:"duration"`
}

type pipelineListResponse struct {
	Items []pipelineResponse `json:"items"`
	Total int                `json:"total"`
}
