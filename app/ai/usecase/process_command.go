package usecase

import (
	"context"
	"fmt"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/app/ai/repository"
	"github.com/KevenAbraham/ai-assistant/app/ai/service"
)

type AIClient interface {
	Complete(ctx context.Context, messages []entity.Message) (string, error)
}

type ProcessCommandInput struct {
	Text      string
	SessionID string
}

type ProcessCommandOutput struct {
	Response string
	Intent   entity.Intent
}

type ProcessCommandUseCase struct {
	conversationRepo repository.ConversationRepository
	memoryRepo       repository.MemoryRepository
	aiClient         AIClient
	contextBuilder   *service.ContextBuilder
}

func NewProcessCommandUseCase(
	convRepo repository.ConversationRepository,
	memRepo repository.MemoryRepository,
	aiClient AIClient,
	cb *service.ContextBuilder,
) *ProcessCommandUseCase {
	return &ProcessCommandUseCase{
		conversationRepo: convRepo,
		memoryRepo:       memRepo,
		aiClient:         aiClient,
		contextBuilder:   cb,
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
	responseText, err := uc.aiClient.Complete(ctx, messages)
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
