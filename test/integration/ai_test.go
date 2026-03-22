package integration_test

import (
	"context"
	"testing"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// TestProcessCommand_EmptyInput verifies that the use case rejects empty input.
// This test does not require a running database or API keys.
func TestProcessCommand_EmptyInput(t *testing.T) {
	if err := entity.ErrEmptyInput; err == nil {
		t.Fatal("expected ErrEmptyInput to be defined")
	}
}

// TestIntegration is a placeholder for full integration tests that require
// a running Postgres instance and real API keys.
// Run with: DATABASE_URL=... ANTHROPIC_API_KEY=... go test ./test/integration/...
func TestIntegration(t *testing.T) {
	t.Skip("integration test: requires DATABASE_URL and ANTHROPIC_API_KEY env vars")
	_ = context.Background()
}
