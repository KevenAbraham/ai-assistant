package service

import (
	"strings"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

type IntentRouter struct{}

func NewIntentRouter() *IntentRouter {
	return &IntentRouter{}
}

func (r *IntentRouter) Route(text string) *entity.Command {
	lower := strings.ToLower(text)
	cmd := &entity.Command{RawText: text}

	switch {
	case isOpenAppIntent(lower):
		cmd.Intent = entity.IntentOpenApp
		cmd.Action = &entity.Action{
			Type:    "open_app",
			Payload: map[string]string{"query": text},
		}
	case isAlarmIntent(lower):
		cmd.Intent = entity.IntentSetAlarm
		cmd.Action = &entity.Action{
			Type:    "set_alarm",
			Payload: map[string]string{"query": text},
		}
	case isSaveMemoryIntent(lower):
		cmd.Intent = entity.IntentSaveMemory
	case isQueryMemoryIntent(lower):
		cmd.Intent = entity.IntentQueryMemory
	default:
		cmd.Intent = entity.IntentChat
	}

	return cmd
}

func isOpenAppIntent(text string) bool {
	return containsAny(text, "abre ", "abrir ", "open ", "launch ")
}

func isAlarmIntent(text string) bool {
	return containsAny(text, "alarme", "lembra", "alarm", "remind")
}

func isSaveMemoryIntent(text string) bool {
	return containsAny(text, "lembra que", "remember that", "save that")
}

func isQueryMemoryIntent(text string) bool {
	return containsAny(text, "o que você sabe", "what do you know", "recall")
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
