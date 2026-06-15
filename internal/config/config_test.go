package config

import (
	"os"
	"testing"
	"time"
)

func clearEnv() {
	os.Unsetenv("TRELLO_API_KEY")
	os.Unsetenv("TRELLO_TOKEN")
	os.Unsetenv("HTTP_TIMEOUT")
}

func setEnv(t *testing.T, key, val string) {
	t.Helper()
	if err := os.Setenv(key, val); err != nil {
		t.Fatal(err)
	}
}

func TestLoad_Success(t *testing.T) {
	clearEnv()
	setEnv(t, "TRELLO_API_KEY", "test-key")
	setEnv(t, "TRELLO_TOKEN", "test-token")
	setEnv(t, "HTTP_TIMEOUT", "15")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.TrelloAPIKey != "test-key" {
		t.Errorf("TrelloAPIKey = %q, want %q", cfg.TrelloAPIKey, "test-key")
	}
	if cfg.TrelloToken != "test-token" {
		t.Errorf("TrelloToken = %q, want %q", cfg.TrelloToken, "test-token")
	}
	if cfg.HTTPTimeout != 15*time.Second {
		t.Errorf("HTTPTimeout = %v, want %v", cfg.HTTPTimeout, 15*time.Second)
	}
}

func TestLoad_DefaultTimeout(t *testing.T) {
	clearEnv()
	setEnv(t, "TRELLO_API_KEY", "key")
	setEnv(t, "TRELLO_TOKEN", "token")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.HTTPTimeout != 30*time.Second {
		t.Errorf("HTTPTimeout = %v, want %v", cfg.HTTPTimeout, 30*time.Second)
	}
}

func TestLoad_MissingAPIKey(t *testing.T) {
	clearEnv()
	setEnv(t, "TRELLO_TOKEN", "token")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing API key")
	}
	if err.Error() != "missing required environment variable: TRELLO_API_KEY" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLoad_MissingToken(t *testing.T) {
	clearEnv()
	setEnv(t, "TRELLO_API_KEY", "key")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for missing token")
	}
	if err.Error() != "missing required environment variable: TRELLO_TOKEN (generate one at https://trello.com/app-key)" {
		t.Errorf("unexpected error message: %v", err)
	}
}

func TestLoad_InvalidTimeout(t *testing.T) {
	clearEnv()
	setEnv(t, "TRELLO_API_KEY", "key")
	setEnv(t, "TRELLO_TOKEN", "token")
	setEnv(t, "HTTP_TIMEOUT", "not-a-number")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid timeout")
	}
}

func TestLoad_NegativeTimeout(t *testing.T) {
	clearEnv()
	setEnv(t, "TRELLO_API_KEY", "key")
	setEnv(t, "TRELLO_TOKEN", "token")
	setEnv(t, "HTTP_TIMEOUT", "-5")

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for negative timeout")
	}
}
