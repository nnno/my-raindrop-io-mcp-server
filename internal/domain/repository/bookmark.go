package repository

import (
	"context"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

//go:generate mockgen -source=bookmark.go -destination=mock/mock_bookmark.go -package=mock
type BookmarkRepository interface {
	Search(ctx context.Context, params entity.SearchParams) ([]entity.Bookmark, error)
	Get(ctx context.Context, id int) (*entity.Bookmark, error)
	Create(ctx context.Context, params entity.CreateParams) (*entity.Bookmark, error)
	Update(ctx context.Context, id int, params entity.UpdateParams) (*entity.Bookmark, error)
	Delete(ctx context.Context, id int) error
}
