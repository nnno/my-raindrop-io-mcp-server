package usecase

import (
	"context"
	"fmt"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/repository"
)

const (
	defaultPerpage = 25
	maxPerpage     = 50
)

var allowedSortValues = map[string]bool{
	"":            true,
	"-created":    true,
	"-lastUpdate": true,
	"title":       true,
	"-title":      true,
	"domain":      true,
}

type BookmarkUsecase struct {
	repo repository.BookmarkRepository
}

func NewBookmarkUsecase(repo repository.BookmarkRepository) *BookmarkUsecase {
	return &BookmarkUsecase{repo: repo}
}

func (u *BookmarkUsecase) Search(ctx context.Context, params entity.SearchParams) ([]entity.Bookmark, error) {
	if params.CollectionID < 0 {
		return nil, entity.NewValidationError(fmt.Sprintf("collection_id must be non-negative, got %d", params.CollectionID))
	}
	if params.Page < 0 {
		return nil, entity.NewValidationError(fmt.Sprintf("page must be non-negative, got %d", params.Page))
	}
	if params.Perpage < 0 || params.Perpage > maxPerpage {
		return nil, entity.NewValidationError(fmt.Sprintf("perpage must be between 0 and %d, got %d", maxPerpage, params.Perpage))
	}
	if params.Perpage == 0 {
		params.Perpage = defaultPerpage
	}
	if !allowedSortValues[params.Sort] {
		return nil, entity.NewValidationError(fmt.Sprintf("invalid sort value: %q", params.Sort))
	}
	return u.repo.Search(ctx, params)
}

func (u *BookmarkUsecase) Get(ctx context.Context, id int) (*entity.Bookmark, error) {
	if id <= 0 {
		return nil, entity.NewValidationError(fmt.Sprintf("id must be positive, got %d", id))
	}
	return u.repo.Get(ctx, id)
}

func (u *BookmarkUsecase) Create(ctx context.Context, params entity.CreateParams) (*entity.Bookmark, error) {
	if params.Link == "" {
		return nil, entity.NewValidationError("link is required")
	}
	return u.repo.Create(ctx, params)
}

func (u *BookmarkUsecase) Update(ctx context.Context, id int, params entity.UpdateParams) (*entity.Bookmark, error) {
	if id <= 0 {
		return nil, entity.NewValidationError(fmt.Sprintf("id must be positive, got %d", id))
	}
	if params.Title == nil && params.Tags == nil && params.CollectionID == nil && params.Note == nil {
		return nil, entity.NewValidationError("at least one field must be specified for update")
	}
	return u.repo.Update(ctx, id, params)
}

func (u *BookmarkUsecase) Delete(ctx context.Context, id int) error {
	if id <= 0 {
		return entity.NewValidationError(fmt.Sprintf("id must be positive, got %d", id))
	}
	return u.repo.Delete(ctx, id)
}
