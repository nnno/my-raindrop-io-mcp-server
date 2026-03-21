package handler

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

func RegisterListCollections(s ToolAdder, uc *usecase.CollectionUsecase) {
	tool := mcp.NewTool("list_collections",
		mcp.WithDescription("Raindropのコレクション一覧を取得する（ルート＋子コレクション）"),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		collections, err := uc.List(ctx)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(map[string]any{"collections": collections})
	})
}

func RegisterCreateCollection(s ToolAdder, uc *usecase.CollectionUsecase) {
	tool := mcp.NewTool("create_collection",
		mcp.WithDescription("新しいコレクションを作成する"),
		mcp.WithString("title", mcp.Required(), mcp.Description("コレクション名")),
		mcp.WithNumber("parent_id", mcp.Description("親コレクションID（省略時はルート）")),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		title, err := req.RequireString("title")
		if err != nil {
			return mcp.NewToolResultError("title is required"), nil
		}
		params := entity.CollectionCreateParams{
			Title:    title,
			ParentID: req.GetInt("parent_id", 0),
		}
		collection, err := uc.Create(ctx, params)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(collection)
	})
}

func RegisterUpdateCollection(s ToolAdder, uc *usecase.CollectionUsecase) {
	tool := mcp.NewTool("update_collection",
		mcp.WithDescription("既存のコレクションを更新する"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("コレクションID")),
		mcp.WithString("title", mcp.Description("新しいコレクション名")),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError("id is required"), nil
		}

		params := entity.CollectionUpdateParams{}
		args := req.GetArguments()
		if v, ok := args["title"]; ok {
			title, ok := v.(string)
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("title must be a string, got %T", v)), nil
			}
			params.Title = &title
		}

		collection, err := uc.Update(ctx, id, params)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(collection)
	})
}
