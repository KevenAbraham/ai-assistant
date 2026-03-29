package service

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// macOSAppNames maps voice-transcription variants to the exact app name
// as it appears in /Applications on macOS.
var macOSAppNames = map[string]string{
	"chrome":        "Google Chrome",
	"crome":         "Google Chrome",
	"google chrome": "Google Chrome",
	"chromium":      "Chromium",
	"firefox":       "Firefox",
	"safari":        "Safari",
	"spotify":       "Spotify",
	"terminal":      "Terminal",
	"finder":        "Finder",
	"code":          "Visual Studio Code",
	"vscode":        "Visual Studio Code",
	"visual studio": "Visual Studio Code",
	"slack":         "Slack",
	"discord":       "Discord",
	"whatsapp":      "WhatsApp",
	"zoom":          "zoom.us",
	"notes":         "Notes",
	"music":         "Music",
}

var linuxAppAliases = map[string]string{
	"chrome":        "google-chrome",
	"crome":         "google-chrome",
	"google chrome": "google-chrome",
}

func normalise(name string) string {
	return strings.ToLower(strings.TrimSpace(name))
}

func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		if w != "" {
			words[i] = strings.ToUpper(w[:1]) + w[1:]
		}
	}
	return strings.Join(words, " ")
}

func resolveMacApp(name string) string {
	if alias, ok := macOSAppNames[normalise(name)]; ok {
		return alias
	}
	return titleCase(normalise(name))
}

func resolveLinuxApp(name string) string {
	if alias, ok := linuxAppAliases[normalise(name)]; ok {
		return alias
	}
	return name
}

func openURLOnDarwin(ctx context.Context, url string) error {
	return exec.CommandContext(ctx, "open", url).Start()
}

func openURLOnLinux(ctx context.Context, url string) error {
	for _, b := range []string{"xdg-open", "sensible-browser", "firefox", "google-chrome", "chromium"} {
		if exec.CommandContext(ctx, b, url).Start() == nil {
			return nil
		}
	}
	return fmt.Errorf("no browser found (tried xdg-open, sensible-browser, firefox, google-chrome, chromium)")
}

func openURLForOS(ctx context.Context, url string) error {
	switch runtime.GOOS {
	case "darwin":
		return openURLOnDarwin(ctx, url)
	default:
		return openURLOnLinux(ctx, url)
	}
}

func openAppOnDarwin(ctx context.Context, name string) error {
	appName := resolveMacApp(name)
	log.Printf("action: open_app %q → open -a %q", name, appName)
	return exec.CommandContext(ctx, "open", "-a", appName).Start()
}

func openAppOnLinux(ctx context.Context, name string) error {
	cmd := resolveLinuxApp(name)
	log.Printf("action: open_app %q → exec %q", name, cmd)
	return exec.CommandContext(ctx, cmd).Start()
}

func openAppForOS(ctx context.Context, name string) error {
	switch runtime.GOOS {
	case "darwin":
		return openAppOnDarwin(ctx, name)
	default:
		return openAppOnLinux(ctx, name)
	}
}

type ActionExecutor struct{}

func NewActionExecutor() *ActionExecutor {
	return &ActionExecutor{}
}

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

func (e *ActionExecutor) openApp(ctx context.Context, name string) (string, error) {
	if err := openAppForOS(ctx, name); err != nil {
		return "", fmt.Errorf("open_app: %w", err)
	}
	return "ok", nil
}

func (e *ActionExecutor) openURL(ctx context.Context, url string) (string, error) {
	log.Printf("action: open_url %q", url)
	if err := openURLForOS(ctx, url); err != nil {
		return "", fmt.Errorf("open_url: %w", err)
	}
	return "ok", nil
}

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
		return nil
	default:
		return fmt.Errorf("unsupported action type: %s", cmd.Action.Type)
	}
}
