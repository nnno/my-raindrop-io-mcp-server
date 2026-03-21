package handler

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
)

func RegisterDelete(s ToolAdder, uc *usecase.BookmarkUsecase) {
	tool := mcp.NewTool("delete_bookmark",
		mcp.WithDescription("ブックマークを削除する"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ブックマークID")),
		mcp.WithDestructiveHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError("id is required"), nil
		}
		if err := uc.Delete(ctx, id); err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultText("deleted"), nil
	})
}
