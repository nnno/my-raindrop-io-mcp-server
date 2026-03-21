package handler

import (
	"context"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

func RegisterSearch(s ToolAdder, uc *usecase.BookmarkUsecase) {
	tool := mcp.NewTool("search_bookmarks",
		mcp.WithDescription("Raindropのブックマークを検索する"),
		mcp.WithNumber("collection_id", mcp.Description("コレクションID（0=全て）")),
		mcp.WithString("search", mcp.Description("検索クエリ")),
		mcp.WithString("sort",
			mcp.Description("ソート順"),
			mcp.Enum("-created", "-lastUpdate", "title", "-title", "domain"),
		),
		mcp.WithNumber("page", mcp.Description("ページ番号（0始まり）")),
		mcp.WithNumber("perpage", mcp.Description("1ページあたりの件数（最大50、デフォルト25）")),
		mcp.WithReadOnlyHintAnnotation(true),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := entity.SearchParams{
			CollectionID: req.GetInt("collection_id", 0),
			Search:       req.GetString("search", ""),
			Sort:         req.GetString("sort", ""),
			Page:         req.GetInt("page", 0),
			Perpage:      req.GetInt("perpage", 0),
		}
		bookmarks, err := uc.Search(ctx, params)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(map[string]any{"bookmarks": bookmarks})
	})
}
