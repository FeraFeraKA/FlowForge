package pipeline

import (
	"context"
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
}

func NewService(repository Repository) *Service {
	return &Service{repository: repository}
}
