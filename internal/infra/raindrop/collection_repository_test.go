package raindrop_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/infra/raindrop"
)

func TestCollectionRepository_List(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)

		w.Header().Set("Content-Type", "application/json")

		switch r.URL.Path {
		case "/collections":
			resp := map[string]any{
				"items": []map[string]any{
					{
						"_id":    1,
						"title":  "Dev",
						"count":  10,
						"parent": map[string]int{"$id": 0},
					},
					{
						"_id":    2,
						"title":  "Design",
						"count":  5,
						"parent": map[string]int{"$id": 0},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		case "/collections/childrens":
			resp := map[string]any{
				"items": []map[string]any{
					{
						"_id":    3,
						"title":  "Go",
						"count":  3,
						"parent": map[string]int{"$id": 1},
					},
				},
			}
			json.NewEncoder(w).Encode(resp)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})

	repo := raindrop.NewCollectionRepository(client)
	collections, err := repo.List(context.Background())

	require.NoError(t, err)
	assert.Len(t, collections, 3)

	assert.Equal(t, "Dev", collections[0].Title)
	assert.Equal(t, 0, collections[0].ParentID)
	assert.Equal(t, "Go", collections[2].Title)
	assert.Equal(t, 1, collections[2].ParentID)
}

func TestCollectionRepository_List_RootError(t *testing.T) {
	client := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error":"server error"}`))
	})

	repo := raindrop.NewCollectionRepository(client)
	_, err := repo.List(context.Background())

	assert.Error(t, err)
}
