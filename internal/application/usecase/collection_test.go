package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/repository/mock"
)

func TestCollectionUsecase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)
	ctx := context.Background()

	expected := []entity.Collection{
		{ID: 1, Title: "Dev", Count: 10, ParentID: 0},
		{ID: 2, Title: "Go", Count: 5, ParentID: 1},
	}
	repo.EXPECT().List(ctx).Return(expected, nil)

	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.List(ctx)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectionUsecase_List_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().List(ctx).Return(nil, errors.New("api error"))

	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.List(ctx)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestCollectionUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)
	ctx := context.Background()

	params := entity.CollectionCreateParams{Title: "New Collection", ParentID: 0}
	expected := &entity.Collection{ID: 10, Title: "New Collection"}
	repo.EXPECT().Create(ctx, params).Return(expected, nil)

	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectionUsecase_Create_EmptyTitle(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)

	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.Create(context.Background(), entity.CollectionCreateParams{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "title is required")
}

func TestCollectionUsecase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)
	ctx := context.Background()

	title := "Updated"
	params := entity.CollectionUpdateParams{Title: &title}
	expected := &entity.Collection{ID: 5, Title: "Updated"}
	repo.EXPECT().Update(ctx, 5, params).Return(expected, nil)

	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.Update(ctx, 5, params)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestCollectionUsecase_Update_NoFields(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)

	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.Update(context.Background(), 1, entity.CollectionUpdateParams{})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "at least one field")
}

func TestCollectionUsecase_Update_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockCollectionRepository(ctrl)

	title := "x"
	uc := usecase.NewCollectionUsecase(repo)
	result, err := uc.Update(context.Background(), 0, entity.CollectionUpdateParams{Title: &title})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "id must be positive")
}
