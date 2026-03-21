package handler

import (
	"errors"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/domain/entity"
)

// ToolAdder abstracts tool registration for testability.
// Satisfied by *server.MCPServer and *mcptest.Server.
type ToolAdder interface {
	AddTool(tool mcp.Tool, handler server.ToolHandlerFunc)
}

func RegisterAll(s ToolAdder, bu *usecase.BookmarkUsecase, cu *usecase.CollectionUsecase) {
	RegisterSearch(s, bu)
	RegisterGet(s, bu)
	RegisterCreate(s, bu)
	RegisterUpdate(s, bu)
	RegisterDelete(s, bu)
	RegisterListCollections(s, cu)
	RegisterCreateCollection(s, cu)
	RegisterUpdateCollection(s, cu)
}

// toolError returns an appropriate error message based on the error kind.
func toolError(err error) *mcp.CallToolResult {
	var domErr *entity.DomainError
	if errors.As(err, &domErr) {
		switch domErr.Kind {
		case entity.ErrValidation:
			return mcp.NewToolResultError(domErr.Message)
		case entity.ErrNotFound:
			return mcp.NewToolResultError("not found: " + domErr.Message)
		case entity.ErrUnauthorized:
			return mcp.NewToolResultError("authentication failed - check RAINDROP_TOKEN")
		case entity.ErrRateLimited:
			return mcp.NewToolResultError("rate limited - please retry after a moment")
		default:
			return mcp.NewToolResultError("internal error: " + domErr.Message)
		}
	}
	return mcp.NewToolResultError(err.Error())
}
