package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/mcptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/repository/mock"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/handler"
)

func setupTestServer(t *testing.T) (*mcptest.Server, *mock.MockBookmarkRepository, *mock.MockCollectionRepository) {
	t.Helper()
	ctrl := gomock.NewController(t)
	bookmarkRepo := mock.NewMockBookmarkRepository(ctrl)
	collectionRepo := mock.NewMockCollectionRepository(ctrl)

	bu := usecase.NewBookmarkUsecase(bookmarkRepo)
	cu := usecase.NewCollectionUsecase(collectionRepo)

	srv := mcptest.NewUnstartedServer(t)
	handler.RegisterAll(srv, bu, cu)

	ctx := context.Background()
	err := srv.Start(ctx)
	require.NoError(t, err)
	t.Cleanup(srv.Close)

	return srv, bookmarkRepo, collectionRepo
}

func TestSearchBookmarks(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Search(gomock.Any(), entity.SearchParams{
		Search:  "golang",
		Perpage: 25,
	}).Return([]entity.Bookmark{
		{ID: 1, Title: "Go Tutorial", Link: "https://go.dev", Tags: []string{"go"}},
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "search_bookmarks",
			Arguments: map[string]any{"search": "golang"},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)

	text := result.Content[0].(mcp.TextContent).Text
	var wrapper struct {
		Bookmarks []entity.Bookmark `json:"bookmarks"`
	}
	err = json.Unmarshal([]byte(text), &wrapper)
	require.NoError(t, err)
	assert.Len(t, wrapper.Bookmarks, 1)
	assert.Equal(t, "Go Tutorial", wrapper.Bookmarks[0].Title)
}

func TestSearchBookmarks_Error(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Search(gomock.Any(), gomock.Any()).Return(nil, errors.New("api error"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "search_bookmarks",
			Arguments: map[string]any{"search": "test"},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "api error")
}

func TestGetBookmark(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Get(gomock.Any(), 42).Return(&entity.Bookmark{
		ID: 42, Title: "Test", Link: "https://example.com",
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_bookmark",
			Arguments: map[string]any{"id": float64(42)},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Test")
}

func TestGetBookmark_Error(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Get(gomock.Any(), 99).Return(nil, errors.New("not found"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_bookmark",
			Arguments: map[string]any{"id": float64(99)},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "not found")
}

func TestCreateBookmark(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Create(gomock.Any(), entity.CreateParams{
		Link:  "https://example.com",
		Title: "Example",
	}).Return(&entity.Bookmark{
		ID: 100, Title: "Example", Link: "https://example.com",
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "create_bookmark",
			Arguments: map[string]any{
				"link":  "https://example.com",
				"title": "Example",
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Example")
}

func TestCreateBookmark_Error(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("create failed"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "create_bookmark",
			Arguments: map[string]any{
				"link": "https://example.com",
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestUpdateBookmark_TitleOnly(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	title := "Updated"
	bookmarkRepo.EXPECT().Update(gomock.Any(), 1, entity.UpdateParams{
		Title: &title,
	}).Return(&entity.Bookmark{
		ID: 1, Title: "Updated",
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_bookmark",
			Arguments: map[string]any{
				"id":    float64(1),
				"title": "Updated",
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Updated")
}

func TestUpdateBookmark_NoteOnly(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	note := "new note"
	bookmarkRepo.EXPECT().Update(gomock.Any(), 2, entity.UpdateParams{
		Note: &note,
	}).Return(&entity.Bookmark{
		ID: 2, Note: "new note",
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_bookmark",
			Arguments: map[string]any{
				"id":   float64(2),
				"note": "new note",
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
}

func TestUpdateBookmark_MultipleFields(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	title := "New Title"
	note := "New Note"
	cid := 5
	tags := []string{"a", "b"}
	bookmarkRepo.EXPECT().Update(gomock.Any(), 3, entity.UpdateParams{
		Title:        &title,
		Note:         &note,
		CollectionID: &cid,
		Tags:         &tags,
	}).Return(&entity.Bookmark{
		ID: 3, Title: "New Title", Note: "New Note", CollectionID: 5, Tags: []string{"a", "b"},
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_bookmark",
			Arguments: map[string]any{
				"id":            float64(3),
				"title":         "New Title",
				"note":          "New Note",
				"collection_id": float64(5),
				"tags":          []any{"a", "b"},
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
}

func TestUpdateBookmark_Error(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("update failed"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_bookmark",
			Arguments: map[string]any{
				"id":    float64(1),
				"title": "x",
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestUpdateBookmark_NoFields(t *testing.T) {
	srv, _, _ := setupTestServer(t)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_bookmark",
			Arguments: map[string]any{
				"id": float64(1),
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "at least one field")
}

func TestListCollections(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	collectionRepo.EXPECT().List(gomock.Any()).Return([]entity.Collection{
		{ID: 1, Title: "Dev", Count: 10},
		{ID: 2, Title: "Go", Count: 5, ParentID: 1},
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "list_collections",
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Dev")
	assert.Contains(t, text, "Go")
}

func TestListCollections_Error(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	collectionRepo.EXPECT().List(gomock.Any()).Return(nil, errors.New("list failed"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "list_collections",
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestCreateCollection(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	collectionRepo.EXPECT().Create(gomock.Any(), entity.CollectionCreateParams{
		Title:    "My Collection",
		ParentID: 0,
	}).Return(&entity.Collection{
		ID: 10, Title: "My Collection",
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "create_collection",
			Arguments: map[string]any{
				"title": "My Collection",
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "My Collection")
}

func TestCreateCollection_WithParent(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	collectionRepo.EXPECT().Create(gomock.Any(), entity.CollectionCreateParams{
		Title:    "Sub Collection",
		ParentID: 5,
	}).Return(&entity.Collection{
		ID: 11, Title: "Sub Collection", ParentID: 5,
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "create_collection",
			Arguments: map[string]any{
				"title":     "Sub Collection",
				"parent_id": float64(5),
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Sub Collection")
}

func TestCreateCollection_Error(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	collectionRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil, errors.New("create failed"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "create_collection",
			Arguments: map[string]any{
				"title": "Test",
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestUpdateCollection(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	title := "Renamed"
	collectionRepo.EXPECT().Update(gomock.Any(), 10, entity.CollectionUpdateParams{
		Title: &title,
	}).Return(&entity.Collection{
		ID: 10, Title: "Renamed",
	}, nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_collection",
			Arguments: map[string]any{
				"id":    float64(10),
				"title": "Renamed",
			},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "Renamed")
}

func TestUpdateCollection_Error(t *testing.T) {
	srv, _, collectionRepo := setupTestServer(t)

	collectionRepo.EXPECT().Update(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("update failed"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_collection",
			Arguments: map[string]any{
				"id":    float64(1),
				"title": "x",
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestUpdateCollection_NoFields(t *testing.T) {
	srv, _, _ := setupTestServer(t)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "update_collection",
			Arguments: map[string]any{
				"id": float64(1),
			},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Contains(t, text, "at least one field")
}

func TestDeleteBookmark(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Delete(gomock.Any(), 42).Return(nil)

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "delete_bookmark",
			Arguments: map[string]any{"id": float64(42)},
		},
	})

	require.NoError(t, err)
	assert.False(t, result.IsError)
	text := result.Content[0].(mcp.TextContent).Text
	assert.Equal(t, "deleted", text)
}

func TestDeleteBookmark_Error(t *testing.T) {
	srv, bookmarkRepo, _ := setupTestServer(t)

	bookmarkRepo.EXPECT().Delete(gomock.Any(), 99).Return(errors.New("not found"))

	result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "delete_bookmark",
			Arguments: map[string]any{"id": float64(99)},
		},
	})

	require.NoError(t, err)
	assert.True(t, result.IsError)
}

func TestToolError_DomainErrorKinds(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		contains string
	}{
		{
			name:     "Unauthorized",
			err:      &entity.DomainError{Kind: entity.ErrUnauthorized, Message: "unauthorized"},
			contains: "authentication failed",
		},
		{
			name:     "RateLimited",
			err:      &entity.DomainError{Kind: entity.ErrRateLimited, Message: "too many requests"},
			contains: "rate limited",
		},
		{
			name:     "Internal",
			err:      &entity.DomainError{Kind: entity.ErrInternal, Message: "server error"},
			contains: "internal error",
		},
		{
			name:     "NotFound",
			err:      &entity.DomainError{Kind: entity.ErrNotFound, Message: "missing"},
			contains: "not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv, bookmarkRepo, _ := setupTestServer(t)

			bookmarkRepo.EXPECT().Get(gomock.Any(), 1).Return(nil, tt.err)

			result, err := srv.Client().CallTool(context.Background(), mcp.CallToolRequest{
				Params: mcp.CallToolParams{
					Name:      "get_bookmark",
					Arguments: map[string]any{"id": float64(1)},
				},
			})

			require.NoError(t, err)
			assert.True(t, result.IsError)
			text := result.Content[0].(mcp.TextContent).Text
			assert.Contains(t, text, tt.contains)
		})
	}
}
