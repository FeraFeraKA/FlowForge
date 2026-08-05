package runmanager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	runmodel "flowforge/backend/internal/run"
)

func (m *Manager) StartRun(ctx context.Context, pipelineID string) (runmodel.Run, error) {
	if err := ctx.Err(); err != nil {
		return runmodel.Run{}, err
	}
	definition, err := m.pipelines.GetByID(ctx, pipelineID)
	if err != nil {
		return runmodel.Run{}, fmt.Errorf("get pipeline %q: %w", pipelineID, err)
	}
	id, err := m.newID()
	if err != nil {
		return runmodel.Run{}, fmt.Errorf("create run ID: %w", err)
	}
	execution := runmodel.New(id, definition, m.now())
	if err := m.runs.Create(ctx, execution); err != nil {
		return runmodel.Run{}, fmt.Errorf("store new run: %w", err)
	}

	runCtx, cancel := context.WithCancel(m.appCtx)
	m.mu.Lock()
	m.active[execution.ID] = cancel
	m.work.Add(1)
	m.mu.Unlock()

	initial := runmodel.Clone(execution)
	go m.execute(runCtx, cancel, &execution)
	return initial, nil
}

func (m *Manager) execute(ctx context.Context, cancel context.CancelFunc, execution *runmodel.Run) {
	defer m.work.Done()
	defer cancel()
	defer func() {
		m.mu.Lock()
		delete(m.active, execution.ID)
		m.mu.Unlock()
	}()

	err := m.scheduler.Execute(ctx, execution)
	saveErr := m.runs.Save(context.WithoutCancel(m.appCtx), *execution)
	if err != nil && !errors.Is(err, context.Canceled) {
		slog.Error("run execution failed", "run_id", execution.ID, "error", err)
	}
	if saveErr != nil {
		slog.Error("save final run snapshot", "run_id", execution.ID, "error", saveErr)
	}
}
