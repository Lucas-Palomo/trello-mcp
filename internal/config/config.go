package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	TrelloAPIKey string
	TrelloToken  string
	HTTPTimeout  time.Duration
}

func Load() (Config, error) {
	apiKey := os.Getenv("TRELLO_API_KEY")
	if apiKey == "" {
		return Config{}, fmt.Errorf("missing required environment variable: TRELLO_API_KEY")
	}

	token := os.Getenv("TRELLO_TOKEN")
	if token == "" {
		return Config{}, fmt.Errorf("missing required environment variable: TRELLO_TOKEN (generate one at https://trello.com/app-key)")
	}

	timeout := 30 * time.Second
	if v := os.Getenv("HTTP_TIMEOUT"); v != "" {
		sec, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid HTTP_TIMEOUT %q: must be an integer seconds", v)
		}
		if sec <= 0 {
			return Config{}, fmt.Errorf("invalid HTTP_TIMEOUT %d: must be positive", sec)
		}
		timeout = time.Duration(sec) * time.Second
	}

	return Config{
		TrelloAPIKey: apiKey,
		TrelloToken:  token,
		HTTPTimeout:  timeout,
	}, nil
}
