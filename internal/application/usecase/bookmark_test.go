package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/repository/mock"
)

func TestBookmarkUsecase_Search(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	expected := []entity.Bookmark{
		{ID: 1, Title: "Go Tutorial", Link: "https://go.dev"},
	}
	// Perpage=0 should be filled with default 25
	repo.EXPECT().Search(ctx, entity.SearchParams{Search: "golang", Perpage: 25}).Return(expected, nil)

	uc := usecase.NewBookmarkUsecase(repo)
	result, err := uc.Search(ctx, entity.SearchParams{Search: "golang"})

	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBookmarkUsecase_Search_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Search(ctx, entity.SearchParams{Search: "test", Perpage: 25}).Return(nil, errors.New("api error"))

	uc := usecase.NewBookmarkUsecase(repo)
	result, err := uc.Search(ctx, entity.SearchParams{Search: "test"})

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBookmarkUsecase_Search_NegativePage(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	_, err := uc.Search(context.Background(), entity.SearchParams{Page: -1})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "page must be non-negative")
}

func TestBookmarkUsecase_Search_PerpageExceedsMax(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	_, err := uc.Search(context.Background(), entity.SearchParams{Perpage: 100})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "perpage must be between")
}

func TestBookmarkUsecase_Search_InvalidSort(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	_, err := uc.Search(context.Background(), entity.SearchParams{Sort: "invalid"})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid sort value")
}

func TestBookmarkUsecase_Search_ValidSorts(t *testing.T) {
	for _, sort := range []string{"", "-created", "-lastUpdate", "title", "-title", "domain"} {
		t.Run("sort="+sort, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			repo := mock.NewMockBookmarkRepository(ctrl)
			ctx := context.Background()

			repo.EXPECT().Search(ctx, entity.SearchParams{Sort: sort, Perpage: 25}).Return(nil, nil)

			uc := usecase.NewBookmarkUsecase(repo)
			_, err := uc.Search(ctx, entity.SearchParams{Sort: sort})
			assert.NoError(t, err)
		})
	}
}

func TestBookmarkUsecase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	expected := &entity.Bookmark{ID: 42, Title: "Test Bookmark"}
	repo.EXPECT().Get(ctx, 42).Return(expected, nil)

	uc := usecase.NewBookmarkUsecase(repo)
	result, err := uc.Get(ctx, 42)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBookmarkUsecase_Get_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)

	_, err := uc.Get(context.Background(), 0)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be positive")

	_, err = uc.Get(context.Background(), -1)
	assert.Error(t, err)
}

func TestBookmarkUsecase_Get_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Get(ctx, 99).Return(nil, errors.New("not found"))

	uc := usecase.NewBookmarkUsecase(repo)
	result, err := uc.Get(ctx, 99)

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestBookmarkUsecase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	params := entity.CreateParams{
		Link:  "https://example.com",
		Title: "Example",
		Tags:  []string{"test"},
	}
	expected := &entity.Bookmark{ID: 1, Title: "Example", Link: "https://example.com"}
	repo.EXPECT().Create(ctx, params).Return(expected, nil)

	uc := usecase.NewBookmarkUsecase(repo)
	result, err := uc.Create(ctx, params)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBookmarkUsecase_Create_EmptyLink(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	_, err := uc.Create(context.Background(), entity.CreateParams{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "link is required")
}

func TestBookmarkUsecase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	title := "Updated Title"
	params := entity.UpdateParams{Title: &title}
	expected := &entity.Bookmark{ID: 1, Title: "Updated Title"}
	repo.EXPECT().Update(ctx, 1, params).Return(expected, nil)

	uc := usecase.NewBookmarkUsecase(repo)
	result, err := uc.Update(ctx, 1, params)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestBookmarkUsecase_Update_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	title := "x"
	_, err := uc.Update(context.Background(), 0, entity.UpdateParams{Title: &title})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be positive")
}

func TestBookmarkUsecase_Update_EmptyParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	_, err := uc.Update(context.Background(), 1, entity.UpdateParams{})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least one field must be specified")
}

func TestBookmarkUsecase_Search_NegativeCollectionID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	_, err := uc.Search(context.Background(), entity.SearchParams{CollectionID: -1})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "collection_id must be non-negative")
}

func TestBookmarkUsecase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, 42).Return(nil)

	uc := usecase.NewBookmarkUsecase(repo)
	err := uc.Delete(ctx, 42)

	assert.NoError(t, err)
}

func TestBookmarkUsecase_Delete_InvalidID(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)

	uc := usecase.NewBookmarkUsecase(repo)
	err := uc.Delete(context.Background(), 0)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "id must be positive")
}

func TestBookmarkUsecase_Delete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	repo := mock.NewMockBookmarkRepository(ctrl)
	ctx := context.Background()

	repo.EXPECT().Delete(ctx, 99).Return(errors.New("not found"))

	uc := usecase.NewBookmarkUsecase(repo)
	err := uc.Delete(ctx, 99)

	assert.Error(t, err)
}
