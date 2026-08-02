package pipeline

import (
	"context"
	"fmt"
	"time"
)

type Repository interface {
	Create(ctx context.Context, pipeline Pipeline) error
	GetByID(ctx context.Context, id string) (Pipeline, error)
	List(ctx context.Context) ([]Pipeline, error)
}

type CreateParams struct {
	Name        string
	MaxParallel int
	Jobs        []JobDefinition
}

type Service struct {
	repository Repository
	newID      func() (string, error)
	now        func() time.Time
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
		newID:      generatePipelineID,
		now:        func() time.Time { return time.Now().UTC() },
	}
}

func (s *Service) Create(ctx context.Context, params CreateParams) (Pipeline, error) {
	if err := ctx.Err(); err != nil {
		return Pipeline{}, err
	}

	pipeline := clonePipeline(Pipeline{
		Name:        params.Name,
		MaxParallel: params.MaxParallel,
		Jobs:        params.Jobs,
	})
	if err := pipeline.Validate(); err != nil {
		return Pipeline{}, err
	}

	id, err := s.newID()
	if err != nil {
		return Pipeline{}, fmt.Errorf("create pipeline ID: %w", err)
	}
	pipeline.ID = id
	pipeline.CreatedAt = s.now()

	if err := s.repository.Create(ctx, pipeline); err != nil {
		return Pipeline{}, fmt.Errorf("store pipeline: %w", err)
	}

	return pipeline, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (Pipeline, error) {
	pipeline, err := s.repository.GetByID(ctx, id)
	if err != nil {
		return Pipeline{}, fmt.Errorf("get pipeline %q: %w", id, err)
	}
	return pipeline, nil
}

func (s *Service) List(ctx context.Context) ([]Pipeline, error) {
	pipelines, err := s.repository.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("list pipelines: %w", err)
	}
	return pipelines, nil
}
