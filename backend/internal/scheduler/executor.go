package scheduler

import (
	"context"
	"errors"

	"flowforge/backend/internal/pipeline"
)

var ErrInvalidExecutorConfig = errors.New("invalid executor config")

type Executor interface {
	Execute(ctx context.Context, definition pipeline.JobDefinition) error
}
