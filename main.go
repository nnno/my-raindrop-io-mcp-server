package main

import (
	"fmt"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/nnno/my-raindrop-io-mcp-server/internal/application/usecase"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/handler"
	"github.com/nnno/my-raindrop-io-mcp-server/internal/infra/raindrop"
)

var version = "dev"

func main() {
	token := os.Getenv("RAINDROP_TOKEN")
	if token == "" {
		fmt.Fprintln(os.Stderr, "RAINDROP_TOKEN environment variable is required")
		os.Exit(1)
	}

	client := raindrop.NewClient(token)
	bookmarkRepo := raindrop.NewBookmarkRepository(client)
	collectionRepo := raindrop.NewCollectionRepository(client)
	bu := usecase.NewBookmarkUsecase(bookmarkRepo)
	cu := usecase.NewCollectionUsecase(collectionRepo)

	s := server.NewMCPServer("raindrop", version)
	handler.RegisterAll(s, bu, cu)

	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
