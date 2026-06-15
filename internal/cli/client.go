package cli

import (
	"fmt"

	"github.com/mark3labs/mcp-go/client"
	"trello-mcp/internal/config"
	"trello-mcp/internal/server"
	"trello-mcp/internal/trello"
)

var newMCPClient = func() (*client.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("config: %w", err)
	}
	tc := trello.NewClient(cfg.TrelloAPIKey, cfg.TrelloToken, cfg.HTTPTimeout)
	srv := server.New(&server.Config{TrelloClient: tc})
	c, err := client.NewInProcessClient(srv.Srv)
	if err != nil {
		return nil, fmt.Errorf("in-process client: %w", err)
	}
	return c, nil
}
