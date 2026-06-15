package main

import (
	"fmt"
	"log"
	"os"

	"trello-mcp/internal/config"
	"trello-mcp/internal/server"
	"trello-mcp/internal/trello"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	client := trello.NewClient(cfg.TrelloAPIKey, cfg.TrelloToken, cfg.HTTPTimeout)

	srv := server.New(&server.Config{
		TrelloClient: client,
	})

	log.Println("trello-mcp server starting")
	if err := srv.Serve(); err != nil {
		fmt.Fprintf(os.Stderr, "server error: %v\n", err)
		os.Exit(1)
	}
}
