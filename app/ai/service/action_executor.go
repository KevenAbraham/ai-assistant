package service

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// appAliases maps common voice-transcription variants and display names to
// the actual executable command. Extend as needed.
var appAliases = map[string]string{
	"chrome":          "google-chrome",
	"crome":           "google-chrome",
	"google chrome":   "google-chrome",
	"google crome":    "google-chrome",
	"chromium":        "chromium",
	"firefox":         "firefox",
	"spotify":         "spotify",
	"nautilus":        "nautilus",
	"files":           "nautilus",
	"terminal":        "gnome-terminal",
	"code":            "code",
	"vscode":          "code",
	"visual studio":   "code",
}

// resolveApp returns the canonical executable name for an app.
func resolveApp(name string) string {
	if alias, ok := appAliases[strings.ToLower(strings.TrimSpace(name))]; ok {
		return alias
	}
	return name
}

// ActionExecutor runs local device actions (open app, set alarm, etc.).
type ActionExecutor struct{}

func NewActionExecutor() *ActionExecutor {
	return &ActionExecutor{}
}

// HandleTool is the entry point called by the LLM tool handler.
func (e *ActionExecutor) HandleTool(ctx context.Context, name string, input map[string]interface{}) (string, error) {
	switch name {
	case "open_app":
		app, _ := input["app_name"].(string)
		if app == "" {
			return "", fmt.Errorf("open_app: missing app_name")
		}
		return e.openApp(ctx, app)
	case "open_url":
		url, _ := input["url"].(string)
		if url == "" {
			return "", fmt.Errorf("open_url: missing url")
		}
		return e.openURL(ctx, url)
	default:
		return "", fmt.Errorf("unknown tool: %s", name)
	}
}

// openApp launches an application by its command name.
// Uses the app alias map to normalise common voice-transcription variants.
func (e *ActionExecutor) openApp(ctx context.Context, app string) (string, error) {
	cmd := resolveApp(app)
	log.Printf("action: open_app %q → exec %q", app, cmd)
	if err := exec.CommandContext(ctx, cmd).Start(); err != nil {
		// Fallback: if the resolved name failed, try the original.
		if cmd != app {
			log.Printf("action: open_app fallback to original %q", app)
			if err2 := exec.CommandContext(ctx, app).Start(); err2 == nil {
				return "ok", nil
			}
		}
		return "", fmt.Errorf("open_app %q: %w", cmd, err)
	}
	return "ok", nil
}

// openURL opens a URL using the first available browser or xdg-open.
func (e *ActionExecutor) openURL(ctx context.Context, url string) (string, error) {
	browsers := []string{"xdg-open", "sensible-browser", "firefox", "google-chrome", "chromium"}
	log.Printf("action: open_url %q", url)
	for _, b := range browsers {
		if err := exec.CommandContext(ctx, b, url).Start(); err == nil {
			return "ok", nil
		}
	}
	return "", fmt.Errorf("open_url: no browser found (tried xdg-open, sensible-browser, firefox, google-chrome, chromium)")
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
		_, err := e.openApp(ctx, app)
		return err
	case "set_alarm":
		// Placeholder: integrate with system notification / cron in production.
		return nil
	default:
		return fmt.Errorf("unsupported action type: %s", cmd.Action.Type)
	}
}
