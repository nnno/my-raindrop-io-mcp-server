package usecase

import (
	"context"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/repository"
)

type CollectionUsecase struct {
	repo repository.CollectionRepository
}

func NewCollectionUsecase(repo repository.CollectionRepository) *CollectionUsecase {
	return &CollectionUsecase{repo: repo}
}

func (u *CollectionUsecase) List(ctx context.Context) ([]entity.Collection, error) {
	return u.repo.List(ctx)
}

func (u *CollectionUsecase) Create(ctx context.Context, params entity.CollectionCreateParams) (*entity.Collection, error) {
	if params.Title == "" {
		return nil, entity.NewValidationError("title is required")
	}
	return u.repo.Create(ctx, params)
}

func (u *CollectionUsecase) Update(ctx context.Context, id int, params entity.CollectionUpdateParams) (*entity.Collection, error) {
	if id <= 0 {
		return nil, entity.NewValidationError("id must be positive")
	}
	if params.Title == nil {
		return nil, entity.NewValidationError("at least one field must be specified for update")
	}
	return u.repo.Update(ctx, id, params)
}
