package scheduler

import (
	"context"
	"fmt"
	"time"

	"flowforge/backend/internal/pipeline"
)

type DelayExecutor struct{}

func (DelayExecutor) Execute(ctx context.Context, definition pipeline.JobDefinition) error {
	if definition.Type != pipeline.ExecutorTypeDelay {
		return fmt.Errorf(
			"%w: delay executor cannot execute type %q",
			ErrInvalidExecutorConfig,
			definition.Type,
		)
	}
	if definition.DelayConfig == nil {
		return fmt.Errorf("%w: delay config is missing", ErrInvalidExecutorConfig)
	}
	if definition.DelayConfig.Duration <= 0 {
		return fmt.Errorf("%w: delay duration must be positive", ErrInvalidExecutorConfig)
	}

	timer := time.NewTimer(definition.DelayConfig.Duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

var _ Executor = DelayExecutor{}
