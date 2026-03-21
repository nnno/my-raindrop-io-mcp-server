package repository

import (
	"context"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

//go:generate mockgen -source=collection.go -destination=mock/mock_collection.go -package=mock
type CollectionRepository interface {
	List(ctx context.Context) ([]entity.Collection, error)
	Create(ctx context.Context, params entity.CollectionCreateParams) (*entity.Collection, error)
	Update(ctx context.Context, id int, params entity.CollectionUpdateParams) (*entity.Collection, error)
}
