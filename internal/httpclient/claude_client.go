package httpclient

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/internal/config"
)

// ClaudeClient wraps the official Anthropic SDK and implements usecase.AIClient.
type ClaudeClient struct {
	client anthropic.Client
	model  anthropic.Model
}

func NewClaudeClient(cfg *config.Config) *ClaudeClient {
	client := anthropic.NewClient(option.WithAPIKey(cfg.AnthropicAPIKey))
	return &ClaudeClient{
		client: client,
		model:  anthropic.Model(cfg.ClaudeModel),
	}
}

func (c *ClaudeClient) Complete(ctx context.Context, messages []entity.Message) (string, error) {
	var systemBlocks []anthropic.TextBlockParam
	var apiMessages []anthropic.MessageParam

	for _, m := range messages {
		switch m.Role {
		case entity.RoleSystem:
			systemBlocks = append(systemBlocks, anthropic.TextBlockParam{Text: m.Content})
		case entity.RoleUser:
			apiMessages = append(apiMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(m.Content)))
		case entity.RoleAssistant:
			apiMessages = append(apiMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(m.Content)))
		}
	}

	params := anthropic.MessageNewParams{
		Model:     c.model,
		MaxTokens: 1024,
		Messages:  apiMessages,
		System:    systemBlocks,
	}

	resp, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return "", fmt.Errorf("claude API: %w", err)
	}

	if len(resp.Content) == 0 {
		return "", entity.ErrAIClientFailure
	}

	return resp.Content[0].AsText().Text, nil
}
