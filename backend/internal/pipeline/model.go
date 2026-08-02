package pipeline

import "time"

type ExecutorType string

const ExecutorTypeDelay ExecutorType = "delay"

type Pipeline struct {
	ID          string
	Name        string
	MaxParallel int
	Jobs        []JobDefinition
	CreatedAt   time.Time
}

type JobDefinition struct {
	ID          string
	Name        string
	Type        ExecutorType
	DependsOn   []string
	Timeout     time.Duration
	MaxRetries  int
	RetryDelay  time.Duration
	DelayConfig *DelayConfig
}

type DelayConfig struct {
	Duration time.Duration
}
