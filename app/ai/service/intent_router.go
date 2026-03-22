package service

import (
	"strings"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// IntentRouter analyses raw text and maps it to a structured Command.
type IntentRouter struct{}

func NewIntentRouter() *IntentRouter {
	return &IntentRouter{}
}

// Route parses raw user text and returns a Command with the detected intent.
// This is a keyword-based implementation; replace with a proper NLU model in production.
func (r *IntentRouter) Route(text string) *entity.Command {
	lower := strings.ToLower(text)

	cmd := &entity.Command{RawText: text}

	switch {
	case strings.Contains(lower, "abre ") || strings.Contains(lower, "abrir ") ||
		strings.Contains(lower, "open ") || strings.Contains(lower, "launch "):
		cmd.Intent = entity.IntentOpenApp
		cmd.Action = &entity.Action{
			Type:    "open_app",
			Payload: map[string]string{"query": text},
		}

	case strings.Contains(lower, "alarme") || strings.Contains(lower, "lembra") ||
		strings.Contains(lower, "alarm") || strings.Contains(lower, "remind"):
		cmd.Intent = entity.IntentSetAlarm
		cmd.Action = &entity.Action{
			Type:    "set_alarm",
			Payload: map[string]string{"query": text},
		}

	case strings.Contains(lower, "lembra que") || strings.Contains(lower, "remember that") ||
		strings.Contains(lower, "save that"):
		cmd.Intent = entity.IntentSaveMemory

	case strings.Contains(lower, "o que você sabe") || strings.Contains(lower, "what do you know") ||
		strings.Contains(lower, "recall"):
		cmd.Intent = entity.IntentQueryMemory

	default:
		cmd.Intent = entity.IntentChat
	}

	return cmd
}
