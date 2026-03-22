package service

import (
	"fmt"
	"strings"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// ContextBuilder assembles the message list sent to the LLM.
// It prepends a system prompt enriched with long-term memories.
type ContextBuilder struct {
	systemPrompt string
}

func NewContextBuilder(systemPrompt string) *ContextBuilder {
	return &ContextBuilder{systemPrompt: systemPrompt}
}

// Build returns the full message list: system prompt + memories + conversation history.
func (b *ContextBuilder) Build(history []entity.Message, memories []*entity.Memory) []entity.Message {
	systemContent := b.systemPrompt
	if len(memories) > 0 {
		var sb strings.Builder
		sb.WriteString("\n\n## Long-term Memory\n")
		for _, m := range memories {
			sb.WriteString(fmt.Sprintf("- %s: %s\n", m.Key, m.Value))
		}
		systemContent += sb.String()
	}

	messages := []entity.Message{
		{Role: entity.RoleSystem, Content: systemContent},
	}
	messages = append(messages, history...)
	return messages
}
