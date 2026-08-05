package runmanager

import (
	"context"
	"fmt"

	runmodel "flowforge/backend/internal/run"
)

func (m *Manager) GetByID(ctx context.Context, id string) (runmodel.Run, error) {
	execution, err := m.runs.GetByID(ctx, id)
	if err != nil {
		return runmodel.Run{}, fmt.Errorf("get run %q: %w", id, err)
	}
	return execution, nil
}

func (m *Manager) List(ctx context.Context) ([]runmodel.Run, error) {
	runs, err := m.runs.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list runs: %w", err)
	}
	return runs, nil
}

func (m *Manager) Cancel(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	m.mu.Lock()
	cancel, active := m.active[id]
	m.mu.Unlock()
	if active {
		cancel()
		return nil
	}

	if _, err := m.runs.GetByID(ctx, id); err != nil {
		return fmt.Errorf("get run to cancel %q: %w", id, err)
	}
	return fmt.Errorf("%w: %s", ErrRunNotActive, id)
}

func (m *Manager) Wait() {
	m.work.Wait()
}
