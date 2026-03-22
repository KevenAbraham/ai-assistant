package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all environment-variable-driven configuration for the application.
// This is the only file in the project that reads os.Getenv.
type Config struct {
	// Database
	DatabaseURL string

	// Anthropic / Claude
	AnthropicAPIKey string
	ClaudeModel     string

	// HTTP server
	HTTPAddr string

	// Application
	SystemPromptPath string

	// Local whisper.cpp server
	WhisperURL string

	// Voice daemon
	RecordSeconds     int
	SilenceThreshold  float64
	SilenceDurationMs int
}

// Load reads config from environment variables and validates required fields.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		AnthropicAPIKey:  os.Getenv("ANTHROPIC_API_KEY"),
		ClaudeModel:      getEnvOr("CLAUDE_MODEL", "claude-haiku-4-5-20251001"),
		HTTPAddr:         getEnvOr("HTTP_ADDR", ":3000"),
		SystemPromptPath: getEnvOr("SYSTEM_PROMPT_PATH", "resources/system_prompt.txt"),
		WhisperURL:        getEnvOr("WHISPER_URL", "http://localhost:9000"),
		RecordSeconds:     getEnvInt("RECORD_SECONDS", 30),
		SilenceThreshold:  getEnvFloat("SILENCE_THRESHOLD", 500.0),
		SilenceDurationMs: getEnvInt("SILENCE_DURATION_MS", 1000),
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *Config) validate() error {
	required := map[string]string{
		"DATABASE_URL":      c.DatabaseURL,
		"ANTHROPIC_API_KEY": c.AnthropicAPIKey,
	}
	for key, val := range required {
		if val == "" {
			return fmt.Errorf("required environment variable %q is not set", key)
		}
	}
	return nil
}

func getEnvOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil && n > 0 {
			return n
		}
	}
	return fallback
}

func getEnvFloat(key string, fallback float64) float64 {
	if v := os.Getenv(key); v != "" {
		if f, err := strconv.ParseFloat(v, 64); err == nil && f >= 0 {
			return f
		}
	}
	return fallback
}
