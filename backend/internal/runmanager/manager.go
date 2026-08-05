package runmanager

import (
	"context"
	"errors"
	"sync"
	"time"

	"flowforge/backend/internal/pipeline"
	runmodel "flowforge/backend/internal/run"
)

var (
	ErrRunNotActive    = errors.New("run is not active")
	ErrRetryNotAllowed = errors.New("run cannot be retried")
)

type PipelineReader interface {
	GetByID(ctx context.Context, id string) (pipeline.Pipeline, error)
}

type RunRepository interface {
	Create(ctx context.Context, execution runmodel.Run) error
	GetByID(ctx context.Context, id string) (runmodel.Run, error)
	List(ctx context.Context) ([]runmodel.Run, error)
}

type Scheduler interface {
	Execute(ctx context.Context, execution *runmodel.Run) error
}

type Manager struct {
	appCtx    context.Context
	pipelines PipelineReader
	runs      RunRepository
	scheduler Scheduler
	newID     func() (string, error)
	now       func() time.Time

	mu     sync.Mutex
	active map[string]context.CancelFunc
	work   sync.WaitGroup
}

func New(
	appCtx context.Context,
	pipelines PipelineReader,
	runs RunRepository,
	scheduler Scheduler,
) *Manager {
	return &Manager{
		appCtx: appCtx, pipelines: pipelines, runs: runs, scheduler: scheduler,
		newID: generateRunID, now: func() time.Time { return time.Now().UTC() },
		active: make(map[string]context.CancelFunc),
	}
}
