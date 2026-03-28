package service

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// ActionExecutor runs local device actions (open app, set alarm, etc.).
type ActionExecutor struct{}

func NewActionExecutor() *ActionExecutor {
	return &ActionExecutor{}
}

// HandleTool is the entry point called by the LLM tool handler.
// It dispatches to the appropriate local action based on the tool name.
func (e *ActionExecutor) HandleTool(ctx context.Context, name string, input map[string]interface{}) (string, error) {
	switch name {
	case "open_app":
		app, _ := input["app_name"].(string)
		if app == "" {
			return "", fmt.Errorf("open_app: missing app_name")
		}
		if err := exec.CommandContext(ctx, "xdg-open", app).Start(); err != nil {
			return "", fmt.Errorf("open_app: %w", err)
		}
		return "ok", nil
	case "open_url":
		url, _ := input["url"].(string)
		if url == "" {
			return "", fmt.Errorf("open_url: missing url")
		}
		if err := exec.CommandContext(ctx, "xdg-open", url).Start(); err != nil {
			return "", fmt.Errorf("open_url: %w", err)
		}
		return "ok", nil
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// Execute is kept for backward compatibility with the existing domain model.
func (e *ActionExecutor) Execute(ctx context.Context, cmd *entity.Command) error {
	if cmd.Action == nil {
		return nil
	}

	switch cmd.Action.Type {
	case "open_app":
		app := cmd.Action.Payload["app"]
		if app == "" {
			return fmt.Errorf("open_app: missing 'app' in payload")
		}
		return exec.CommandContext(ctx, "xdg-open", app).Start()
	case "set_alarm":
		// Placeholder: integrate with system notification / cron in production.
		return nil
	default:
		return fmt.Errorf("unsupported action type: %s", cmd.Action.Type)
	}
}
