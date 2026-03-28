package usecase

import (
	"context"
	"fmt"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/app/ai/repository"
	"github.com/KevenAbraham/ai-assistant/app/ai/service"
)

// ToolHandler executes a named tool call and returns a result string.
type ToolHandler func(ctx context.Context, name string, input map[string]interface{}) (string, error)

type AIClient interface {
	Complete(ctx context.Context, messages []entity.Message) (string, error)
	CompleteWithTools(ctx context.Context, messages []entity.Message, tools []entity.Tool, handler ToolHandler) (string, error)
}

type ProcessCommandInput struct {
	Text      string
	SessionID string
}

type ProcessCommandOutput struct {
	Response string
	Intent   entity.Intent
}

// availableTools lists the tools exposed to the LLM.
var availableTools = []entity.Tool{
	{
		Name:        "open_app",
		Description: "Abre um aplicativo instalado no computador do usuário.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"app_name": map[string]interface{}{
					"type":        "string",
					"description": "Nome ou comando do aplicativo (ex: firefox, spotify, nautilus)",
				},
			},
			"required": []string{"app_name"},
		},
	},
	{
		Name:        "open_url",
		Description: "Abre uma URL no navegador padrão do usuário.",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"url": map[string]interface{}{
					"type":        "string",
					"description": "URL completa a abrir (ex: https://youtube.com)",
				},
			},
			"required": []string{"url"},
		},
	},
}

type ProcessCommandUseCase struct {
	conversationRepo repository.ConversationRepository
	memoryRepo       repository.MemoryRepository
	aiClient         AIClient
	contextBuilder   *service.ContextBuilder
	actionExecutor   *service.ActionExecutor
}

func NewProcessCommandUseCase(
	convRepo repository.ConversationRepository,
	memRepo repository.MemoryRepository,
	aiClient AIClient,
	cb *service.ContextBuilder,
	executor *service.ActionExecutor,
) *ProcessCommandUseCase {
	return &ProcessCommandUseCase{
		conversationRepo: convRepo,
		memoryRepo:       memRepo,
		aiClient:         aiClient,
		contextBuilder:   cb,
		actionExecutor:   executor,
	}
}

func (uc *ProcessCommandUseCase) Execute(ctx context.Context, input ProcessCommandInput) (*ProcessCommandOutput, error) {
	if input.Text == "" {
		return nil, entity.ErrEmptyInput
	}

	conv, err := uc.conversationRepo.FindBySessionID(ctx, input.SessionID)
	if err != nil {
		conv = &entity.Conversation{
			SessionID: input.SessionID,
			Messages:  []entity.Message{},
		}
	}

	userMsg := entity.Message{
		Role:    entity.RoleUser,
		Content: input.Text,
	}
	conv.Messages = append(conv.Messages, userMsg)

	memories, _ := uc.memoryRepo.FindAll(ctx)
	messages := uc.contextBuilder.Build(conv.Messages, memories)

	handler := func(ctx context.Context, name string, toolInput map[string]interface{}) (string, error) {
		return uc.actionExecutor.HandleTool(ctx, name, toolInput)
	}

	responseText, err := uc.aiClient.CompleteWithTools(ctx, messages, availableTools, handler)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", entity.ErrAIClientFailure, err)
	}

	assistantMsg := entity.Message{
		Role:    entity.RoleAssistant,
		Content: responseText,
	}
	conv.Messages = append(conv.Messages, assistantMsg)

	if saveErr := uc.conversationRepo.Save(ctx, conv); saveErr != nil {
		_ = saveErr
	}

	return &ProcessCommandOutput{
		Response: responseText,
		Intent:   entity.IntentChat,
	}, nil
}
