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

func (e *ActionExecutor) Execute(ctx context.Context, cmd *entity.Command) error {
	if cmd.Action == nil {
		return nil
	}

	switch cmd.Action.Type {
	case "open_app":
		return e.openApp(ctx, cmd.Action.Payload)
	case "set_alarm":
		return e.setAlarm(ctx, cmd.Action.Payload)
	default:
		return fmt.Errorf("unsupported action type: %s", cmd.Action.Type)
	}
}

func (e *ActionExecutor) openApp(_ context.Context, payload map[string]string) error {
	app, ok := payload["app"]
	if !ok || app == "" {
		return fmt.Errorf("open_app: missing 'app' in payload")
	}
	return exec.Command("xdg-open", app).Start()
}

func (e *ActionExecutor) setAlarm(_ context.Context, payload map[string]string) error {
	// Placeholder: integrate with system notification / cron in production.
	_ = payload
	return nil
}
