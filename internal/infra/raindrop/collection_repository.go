package raindrop

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"slices"

	"golang.org/x/sync/errgroup"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

type CollectionRepository struct {
	client *Client
}

func NewCollectionRepository(client *Client) *CollectionRepository {
	return &CollectionRepository{client: client}
}

type apiCollection struct {
	ID     int              `json:"_id"`
	Title  string           `json:"title"`
	Count  int              `json:"count"`
	Parent apiCollectionRef `json:"parent"`
}

type apiCollectionsResponse struct {
	Items []apiCollection `json:"items"`
}

type apiCollectionSingleResponse struct {
	Item apiCollection `json:"item"`
}

func (ac *apiCollection) toEntity() entity.Collection {
	return entity.Collection{
		ID:       ac.ID,
		Title:    ac.Title,
		Count:    ac.Count,
		ParentID: ac.Parent.ID,
	}
}

func (r *CollectionRepository) List(ctx context.Context) ([]entity.Collection, error) {
	var rootResp, childResp apiCollectionsResponse

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		data, err := r.client.doRequest(ctx, "GET", "/collections", nil)
		if err != nil {
			return fmt.Errorf("fetching root collections: %w", err)
		}
		if err := json.Unmarshal(data, &rootResp); err != nil {
			return fmt.Errorf("decoding root collections: %w", err)
		}
		return nil
	})

	g.Go(func() error {
		data, err := r.client.doRequest(ctx, "GET", "/collections/childrens", nil)
		if err != nil {
			return fmt.Errorf("fetching child collections: %w", err)
		}
		if err := json.Unmarshal(data, &childResp); err != nil {
			return fmt.Errorf("decoding child collections: %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		return nil, err
	}

	all := slices.Concat(rootResp.Items, childResp.Items)
	collections := make([]entity.Collection, len(all))
	for i, ac := range all {
		collections[i] = ac.toEntity()
	}
	return collections, nil
}

func (r *CollectionRepository) Create(ctx context.Context, params entity.CollectionCreateParams) (*entity.Collection, error) {
	body := map[string]any{
		"title": params.Title,
	}
	if params.ParentID != 0 {
		body["parent"] = map[string]int{"$id": params.ParentID}
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding create body: %w", err)
	}

	data, err := r.client.doRequest(ctx, "POST", "/collection", bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var resp apiCollectionSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decoding create response: %w", err)
	}

	col := resp.Item.toEntity()
	return &col, nil
}

func (r *CollectionRepository) Update(ctx context.Context, id int, params entity.CollectionUpdateParams) (*entity.Collection, error) {
	body := map[string]any{}
	if params.Title != nil {
		body["title"] = *params.Title
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("encoding update body: %w", err)
	}

	path := fmt.Sprintf("/collection/%d", id)
	data, err := r.client.doRequest(ctx, "PUT", path, bytes.NewReader(jsonBody))
	if err != nil {
		return nil, err
	}

	var resp apiCollectionSingleResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("decoding update response: %w", err)
	}

	col := resp.Item.toEntity()
	return &col, nil
}
