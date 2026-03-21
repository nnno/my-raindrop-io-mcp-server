package raindrop_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/infra/raindrop"
)

func newTestClient(t *testing.T, handler http.HandlerFunc) *raindrop.Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return raindrop.NewClient("test-token", raindrop.WithBaseURL(srv.URL))
}

func TestBookmarkRepository_Search(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Contains(t, r.URL.Path, "/raindrops/0")
		assert.Equal(t, "golang", r.URL.Query().Get("search"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		resp := map[string]any{
			"items": []map[string]any{
				{
					"_id":        1,
					"title":      "Go Tutorial",
					"link":       "https://go.dev",
					"excerpt":    "Learn Go",
					"note":       "",
					"tags":       []string{"go"},
					"type":       "link",
					"created":    "2024-01-01T00:00:00.000Z",
					"lastUpdate": "2024-01-01T00:00:00.000Z",
					"collection": map[string]int{"$id": 0},
					"highlights": []map[string]string{},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	repo := raindrop.NewBookmarkRepository(client)
	bookmarks, err := repo.Search(context.Background(), entity.SearchParams{Search: "golang"})

	require.NoError(t, err)
	assert.Len(t, bookmarks, 1)
	assert.Equal(t, "Go Tutorial", bookmarks[0].Title)
	assert.Equal(t, "https://go.dev", bookmarks[0].Link)
	assert.Equal(t, []string{"go"}, bookmarks[0].Tags)
	assert.Equal(t, entity.BookmarkTypeLink, bookmarks[0].Type)
	assert.Equal(t, time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC), bookmarks[0].Created)
}

func TestBookmarkRepository_Get(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/raindrop/42", r.URL.Path)

		resp := map[string]any{
			"item": map[string]any{
				"_id":        42,
				"title":      "Test Bookmark",
				"link":       "https://example.com",
				"excerpt":    "",
				"note":       "my note",
				"tags":       []string{},
				"type":       "article",
				"created":    "2024-01-01T00:00:00.000Z",
				"lastUpdate": "2024-01-01T00:00:00.000Z",
				"collection": map[string]int{"$id": 5},
				"highlights": []map[string]string{
					{"text": "important", "note": "highlighted"},
				},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	repo := raindrop.NewBookmarkRepository(client)
	bookmark, err := repo.Get(context.Background(), 42)

	require.NoError(t, err)
	assert.Equal(t, 42, bookmark.ID)
	assert.Equal(t, "Test Bookmark", bookmark.Title)
	assert.Equal(t, "my note", bookmark.Note)
	assert.Equal(t, 5, bookmark.CollectionID)
	assert.Equal(t, entity.BookmarkTypeArticle, bookmark.Type)
	assert.Len(t, bookmark.Highlights, 1)
	assert.Equal(t, "important", bookmark.Highlights[0].Text)
}

func TestBookmarkRepository_Create(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "/raindrop", r.URL.Path)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]any
		json.Unmarshal(body, &reqBody)
		assert.Equal(t, "https://example.com", reqBody["link"])
		assert.Equal(t, "Example", reqBody["title"])
		assert.Equal(t, []any{"test"}, reqBody["tags"])

		resp := map[string]any{
			"item": map[string]any{
				"_id":        100,
				"title":      "Example",
				"link":       "https://example.com",
				"excerpt":    "",
				"note":       "",
				"tags":       []string{"test"},
				"type":       "link",
				"created":    "2024-01-01T00:00:00.000Z",
				"lastUpdate": "2024-01-01T00:00:00.000Z",
				"collection": map[string]int{"$id": 0},
				"highlights": []map[string]string{},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	repo := raindrop.NewBookmarkRepository(client)
	bookmark, err := repo.Create(context.Background(), entity.CreateParams{
		Link:  "https://example.com",
		Title: "Example",
		Tags:  []string{"test"},
	})

	require.NoError(t, err)
	assert.Equal(t, 100, bookmark.ID)
	assert.Equal(t, "Example", bookmark.Title)
}

func TestBookmarkRepository_CreateWithCollection(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]any
		json.Unmarshal(body, &reqBody)
		assert.Equal(t, "https://example.com", reqBody["link"])
		col := reqBody["collection"].(map[string]any)
		assert.Equal(t, float64(10), col["$id"])

		resp := map[string]any{
			"item": map[string]any{
				"_id":        101,
				"title":      "",
				"link":       "https://example.com",
				"excerpt":    "",
				"note":       "",
				"tags":       []string{},
				"type":       "link",
				"created":    "2024-01-01T00:00:00.000Z",
				"lastUpdate": "2024-01-01T00:00:00.000Z",
				"collection": map[string]int{"$id": 10},
				"highlights": []map[string]string{},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	repo := raindrop.NewBookmarkRepository(client)
	bookmark, err := repo.Create(context.Background(), entity.CreateParams{
		Link:         "https://example.com",
		CollectionID: 10,
	})

	require.NoError(t, err)
	assert.Equal(t, 10, bookmark.CollectionID)
}

func TestBookmarkRepository_Update(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "PUT", r.Method)
		assert.Equal(t, "/raindrop/100", r.URL.Path)

		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]any
		json.Unmarshal(body, &reqBody)
		assert.Equal(t, "Updated", reqBody["title"])

		resp := map[string]any{
			"item": map[string]any{
				"_id":        100,
				"title":      "Updated",
				"link":       "https://example.com",
				"excerpt":    "",
				"note":       "",
				"tags":       []string{},
				"type":       "link",
				"created":    "2024-01-01T00:00:00.000Z",
				"lastUpdate": "2024-01-02T00:00:00.000Z",
				"collection": map[string]int{"$id": 0},
				"highlights": []map[string]string{},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	repo := raindrop.NewBookmarkRepository(client)
	title := "Updated"
	bookmark, err := repo.Update(context.Background(), 100, entity.UpdateParams{Title: &title})

	require.NoError(t, err)
	assert.Equal(t, "Updated", bookmark.Title)
}

func TestBookmarkRepository_UpdateWithTagsAndCollection(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var reqBody map[string]any
		json.Unmarshal(body, &reqBody)
		assert.Equal(t, []any{"go", "mcp"}, reqBody["tags"])
		col := reqBody["collection"].(map[string]any)
		assert.Equal(t, float64(7), col["$id"])

		resp := map[string]any{
			"item": map[string]any{
				"_id":        100,
				"title":      "Test",
				"link":       "https://example.com",
				"excerpt":    "",
				"note":       "",
				"tags":       []string{"go", "mcp"},
				"type":       "link",
				"created":    "2024-01-01T00:00:00.000Z",
				"lastUpdate": "2024-01-02T00:00:00.000Z",
				"collection": map[string]int{"$id": 7},
				"highlights": []map[string]string{},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	})

	repo := raindrop.NewBookmarkRepository(client)
	tags := []string{"go", "mcp"}
	cid := 7
	bookmark, err := repo.Update(context.Background(), 100, entity.UpdateParams{
		Tags:         &tags,
		CollectionID: &cid,
	})

	require.NoError(t, err)
	assert.Equal(t, []string{"go", "mcp"}, bookmark.Tags)
	assert.Equal(t, 7, bookmark.CollectionID)
}

func TestBookmarkRepository_APIError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"errorMessage":"Unauthorized"}`))
	})

	repo := raindrop.NewBookmarkRepository(client)
	_, err := repo.Search(context.Background(), entity.SearchParams{})

	assert.Error(t, err)
	var domErr *entity.DomainError
	require.ErrorAs(t, err, &domErr)
	assert.Equal(t, entity.ErrUnauthorized, domErr.Kind)
}

func TestBookmarkRepository_Delete(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "/raindrop/42", r.URL.Path)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result":true}`))
	})

	repo := raindrop.NewBookmarkRepository(client)
	err := repo.Delete(context.Background(), 42)

	require.NoError(t, err)
}

func TestBookmarkRepository_Delete_Error(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"errorMessage":"not found"}`))
	})

	repo := raindrop.NewBookmarkRepository(client)
	err := repo.Delete(context.Background(), 999)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "404")
}
