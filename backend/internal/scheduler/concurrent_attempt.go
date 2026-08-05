package scheduler

import (
	"context"
	"fmt"
	"sync"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

type attemptMessage struct {
	jobID string
	err   error
}

type retryMessage struct {
	jobID string
	err   error
}

type concurrentRuntime struct {
	ctx          context.Context
	cancel       context.CancelFunc
	state        executionState
	results      chan attemptMessage
	retries      chan retryMessage
	workers      sync.WaitGroup
	active       int
	retryWaiting int
}

func (s *ConcurrentScheduler) launchAttempt(
	runtime *concurrentRuntime,
	job *runmodel.JobRun,
	definition pipeline.JobDefinition,
) error {
	executor, err := s.base.executorFor(definition.Type)
	if err != nil {
		return err
	}
	if err := job.TransitionTo(runmodel.JobStatusRunning); err != nil {
		return fmt.Errorf("start job %q: %w", job.JobID, err)
	}
	job.Attempt++
	job.Progress = 0
	job.Error = ""
	if job.StartedAt == nil {
		startedAt := s.base.now()
		job.StartedAt = &startedAt
	}
	if err := s.notify(runtime.state.execution); err != nil {
		return fmt.Errorf("observe started job %q: %w", job.JobID, err)
	}

	runtime.active++
	runtime.workers.Add(1)
	jobID := job.JobID
	go func() {
		defer runtime.workers.Done()
		attemptCtx := runtime.ctx
		var cancel context.CancelFunc
		if definition.Timeout > 0 {
			attemptCtx, cancel = context.WithTimeout(runtime.ctx, definition.Timeout)
		}
		executionErr := executor.Execute(attemptCtx, definition)
		if cancel != nil {
			cancel()
		}
		runtime.results <- attemptMessage{jobID: jobID, err: executionErr}
	}()
	return nil
}
