package handler

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
)

func RegisterGet(s ToolAdder, uc *usecase.BookmarkUsecase) {
	tool := mcp.NewTool("get_bookmark",
		mcp.WithDescription("IDを指定してブックマークの詳細を取得する"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ブックマークID")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError("id is required"), nil
		}
		bookmark, err := uc.Get(ctx, id)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(bookmark)
	})
}
