package config

import (
	"os"
	"testing"
)

func TestGetEnvReturnsDefaultWhenMissing(t *testing.T) {
	os.Unsetenv("SERVER_PORT")
	got := getEnv("SERVER_PORT", "8080")
	if got != "8080" {
		t.Fatalf("expected default 8080, got %q", got)
	}
}

func TestGetEnvReturnsEnvValue(t *testing.T) {
	os.Setenv("SERVER_PORT", "9090")
	got := getEnv("SERVER_PORT", "8080")
	if got != "9090" {
		t.Fatalf("expected 9090, got %q", got)
	}
}
