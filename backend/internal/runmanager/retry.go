package runmanager

import (
	"context"
	"fmt"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

func (m *Manager) Retry(ctx context.Context, id string) (runmodel.Run, error) {
	source, err := m.runs.GetByID(ctx, id)
	if err != nil {
		return runmodel.Run{}, fmt.Errorf("get run to retry %q: %w", id, err)
	}
	if source.Status != runmodel.RunStatusFailed && source.Status != runmodel.RunStatusCancelled {
		return runmodel.Run{}, fmt.Errorf("%w: run %s is %s", ErrRetryNotAllowed, id, source.Status)
	}

	definition := pipeline.Pipeline{
		ID:          source.PipelineID,
		Name:        source.PipelineName,
		MaxParallel: source.MaxParallel,
		Jobs:        source.JobDefinitions,
	}
	return m.start(ctx, definition)
}
