package raindrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strconv"
	"time"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

type BookmarkRepository struct {
	client *Client
}

func NewBookmarkRepository(client *Client) *BookmarkRepository {
	return &BookmarkRepository{client: client}
}

// API response types

type apiHighlight struct {
	Text string `json:"text"`
	Note string `json:"note"`
}

type apiBookmark struct {
	ID         int              `json:"_id"`
	Title      string           `json:"title"`
	Link       string           `json:"link"`
	Excerpt    string           `json:"excerpt"`
	Note       string           `json:"note"`
	Tags       []string         `json:"tags"`
	Type       string           `json:"type"`
	Created    string           `json:"created"`
	LastUpdate string           `json:"lastUpdate"`
	Collection apiCollectionRef `json:"collection"`
	Highlights []apiHighlight   `json:"highlights"`
}

type apiSearchResponse struct {
	Items []apiBookmark `json:"items"`
}

type apiSingleResponse struct {
	Item apiBookmark `json:"item"`
}

func parseTime(s string) time.Time {
	t, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		slog.Warn("failed to parse time", "value", s, "error", err)
	}
	return t
}

func toHighlights(hs []apiHighlight) []entity.Highlight {
	highlights := make([]entity.Highlight, len(hs))
	for i, h := range hs {
		highlights[i] = entity.Highlight{Text: h.Text, Note: h.Note}
	}
	return highlights
}

func (ab *apiBookmark) toEntity() entity.Bookmark {
	return entity.Bookmark{
		ID:           ab.ID,
		Title:        ab.Title,
		Link:         ab.Link,
		Excerpt:      ab.Excerpt,
		Note:         ab.Note,
		Tags:         ab.Tags,
		Type:         entity.BookmarkType(ab.Type),
		Created:      parseTime(ab.Created),
		LastUpdate:   parseTime(ab.LastUpdate),
		CollectionID: ab.Collection.ID,
		Highlights: toHighlights(ab.Highlights),
	}
}

func (r *BookmarkRepository) Search(ctx context.Context, params entity.SearchParams) ([]entity.Bookmark, error) {
	q := url.Values{}
	if params.Search != "" {
		q.Set("search", params.Search)
	}
	if params.Sort != "" {
		q.Set("sort", params.Sort)
	}
	if params.Page > 0 {
		q.Set("page", strconv.Itoa(params.Page))
	}
	if params.Perpage > 0 {
		q.Set("perpage", strconv.Itoa(params.Perpage))
	}

	path := fmt.Sprintf("/raindrops/%d", params.CollectionID)
	if len(q) > 0 {
		path += "?" + q.Encode()
	}

	data, err := r.client.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp apiSearchResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decoding search response: %w", err)
	}

	bookmarks := make([]entity.Bookmark, len(resp.Items))
	for i, ab := range resp.Items {
		bookmarks[i] = ab.toEntity()
	}
	return bookmarks, nil
}

func (r *BookmarkRepository) Get(ctx context.Context, id int) (*entity.Bookmark, error) {
	path := fmt.Sprintf("/raindrop/%d", id)

	data, err := r.client.doRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, err
	}

	var resp apiSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decoding get response: %w", err)
	}

	bookmark := resp.Item.toEntity()
	return &bookmark, nil
}

func (r *BookmarkRepository) Create(ctx context.Context, params entity.CreateParams) (*entity.Bookmark, error) {
	body := map[string]any{
		"link": params.Link,
	}
	if params.Title != "" {
		body["title"] = params.Title
	}
	if params.Tags != nil {
		body["tags"] = params.Tags
	}
	if params.CollectionID != 0 {
		body["collection"] = map[string]int{"$id": params.CollectionID}
	}
	if params.Note != "" {
		body["note"] = params.Note
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding create body: %w", err)
	}

	data, err := r.client.doRequest(ctx, "POST", "/raindrop", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var resp apiSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decoding create response: %w", err)
	}

	bookmark := resp.Item.toEntity()
	return &bookmark, nil
}

func (r *BookmarkRepository) Update(ctx context.Context, id int, params entity.UpdateParams) (*entity.Bookmark, error) {
	body := map[string]any{}
	if params.Title != nil {
		body["title"] = *params.Title
	}
	if params.Tags != nil {
		body["tags"] = *params.Tags
	}
	if params.CollectionID != nil {
		body["collection"] = map[string]int{"$id": *params.CollectionID}
	}
	if params.Note != nil {
		body["note"] = *params.Note
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding update body: %w", err)
	}

	path := fmt.Sprintf("/raindrop/%d", id)
	data, err := r.client.doRequest(ctx, "PUT", path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var resp apiSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decoding update response: %w", err)
	}

	bookmark := resp.Item.toEntity()
	return &bookmark, nil
}

func (r *BookmarkRepository) Delete(ctx context.Context, id int) error {
	path := fmt.Sprintf("/raindrop/%d", id)
	_, err := r.client.doRequest(ctx, "DELETE", path, nil)
	return err
}
