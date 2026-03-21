package handler

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

func RegisterCreate(s ToolAdder, uc *usecase.BookmarkUsecase) {
	tool := mcp.NewTool("create_bookmark",
		mcp.WithDescription("新しいブックマークを作成する"),
		mcp.WithString("link", mcp.Required(), mcp.Description("ブックマークURL")),
		mcp.WithString("title", mcp.Description("タイトル")),
		mcp.WithArray("tags", mcp.Description("タグの配列"), mcp.WithStringItems()),
		mcp.WithNumber("collection_id", mcp.Description("コレクションID")),
		mcp.WithString("note", mcp.Description("メモ")),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		link, err := req.RequireString("link")
		if err != nil {
			return mcp.NewToolResultError("link is required"), nil
		}
		params := entity.CreateParams{
			Link:         link,
			Title:        req.GetString("title", ""),
			Tags:         req.GetStringSlice("tags", nil),
			CollectionID: req.GetInt("collection_id", 0),
			Note:         req.GetString("note", ""),
		}
		bookmark, err := uc.Create(ctx, params)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(bookmark)
	})
}
