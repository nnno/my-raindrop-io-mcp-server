package handler

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

func RegisterUpdate(s ToolAdder, uc *usecase.BookmarkUsecase) {
	tool := mcp.NewTool("update_bookmark",
		mcp.WithDescription("既存のブックマークを更新する"),
		mcp.WithNumber("id", mcp.Required(), mcp.Description("ブックマークID")),
		mcp.WithString("title", mcp.Description("新しいタイトル")),
		mcp.WithArray("tags", mcp.Description("新しいタグの配列"), mcp.WithStringItems()),
		mcp.WithNumber("collection_id", mcp.Description("移動先コレクションID")),
		mcp.WithString("note", mcp.Description("新しいメモ")),
	)
	s.AddTool(tool, func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id, err := req.RequireInt("id")
		if err != nil {
			return mcp.NewToolResultError("id is required"), nil
		}

		// GetString/GetInt cannot distinguish "not provided" from "zero value",
		// so we use GetArguments() to detect which fields were explicitly set.
		params := entity.UpdateParams{}
		args := req.GetArguments()
		if v, ok := args["title"]; ok {
			title, ok := v.(string)
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("title must be a string, got %T", v)), nil
			}
			params.Title = &title
		}
		if _, ok := args["tags"]; ok {
			tags := req.GetStringSlice("tags", nil)
			params.Tags = &tags
		}
		if v, ok := args["collection_id"]; ok {
			f, ok := v.(float64)
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("collection_id must be a number, got %T", v)), nil
			}
			cid := int(f)
			params.CollectionID = &cid
		}
		if v, ok := args["note"]; ok {
			note, ok := v.(string)
			if !ok {
				return mcp.NewToolResultError(fmt.Sprintf("note must be a string, got %T", v)), nil
			}
			params.Note = &note
		}

		bookmark, err := uc.Update(ctx, id, params)
		if err != nil {
			return toolError(err), nil
		}
		return mcp.NewToolResultJSON(bookmark)
	})
}
