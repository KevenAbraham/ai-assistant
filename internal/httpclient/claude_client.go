package httpclient

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/app/ai/usecase"
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

// parseMessages separates system blocks from the conversation messages.
func (c *ClaudeClient) parseMessages(messages []entity.Message) ([]anthropic.TextBlockParam, []anthropic.MessageParam) {
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
	return systemBlocks, apiMessages
}

func (c *ClaudeClient) Complete(ctx context.Context, messages []entity.Message) (string, error) {
	systemBlocks, apiMessages := c.parseMessages(messages)

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

// CompleteWithTools runs a multi-turn loop that lets Claude invoke local tools
// until it produces a final text response.
func (c *ClaudeClient) CompleteWithTools(
	ctx context.Context,
	messages []entity.Message,
	tools []entity.Tool,
	handler usecase.ToolHandler,
) (string, error) {
	systemBlocks, apiMessages := c.parseMessages(messages)
	apiTools := buildToolParams(tools)

	for {
		params := anthropic.MessageNewParams{
			Model:     c.model,
			MaxTokens: 1024,
			Messages:  apiMessages,
			System:    systemBlocks,
			Tools:     apiTools,
		}

		message, err := c.client.Messages.New(ctx, params)
		if err != nil {
			return "", fmt.Errorf("claude API: %w", err)
		}

		// Append the assistant turn unconditionally so it's in context for the next call.
		apiMessages = append(apiMessages, message.ToParam())

		// No tool calls → Claude produced the final text response.
		if message.StopReason != anthropic.StopReasonToolUse {
			for _, block := range message.Content {
				if tb, ok := block.AsAny().(anthropic.TextBlock); ok {
					return tb.Text, nil
				}
			}
			return "", entity.ErrAIClientFailure
		}

		// Execute each tool call and collect results.
		var toolResults []anthropic.ContentBlockParamUnion
		for _, block := range message.Content {
			variant, ok := block.AsAny().(anthropic.ToolUseBlock)
			if !ok {
				continue
			}

			var input map[string]interface{}
			if err := json.Unmarshal([]byte(variant.JSON.Input.Raw()), &input); err != nil {
				input = map[string]interface{}{}
			}

			result, toolErr := handler(ctx, variant.Name, input)
			isError := toolErr != nil
			if isError {
				result = toolErr.Error()
			}
			toolResults = append(toolResults, anthropic.NewToolResultBlock(variant.ID, result, isError))
		}

		apiMessages = append(apiMessages, anthropic.NewUserMessage(toolResults...))
	}
}

// buildToolParams converts entity.Tool slice to the SDK's ToolUnionParam slice.
func buildToolParams(tools []entity.Tool) []anthropic.ToolUnionParam {
	params := make([]anthropic.ToolUnionParam, len(tools))
	for i, t := range tools {
		properties, _ := t.Parameters["properties"]
		required, _ := t.Parameters["required"].([]string)

		tool := anthropic.ToolParam{
			Name:        t.Name,
			Description: anthropic.String(t.Description),
			InputSchema: anthropic.ToolInputSchemaParam{
				Properties: properties,
				Required:   required,
			},
		}
		params[i] = anthropic.ToolUnionParam{OfTool: &tool}
	}
	return params
}
