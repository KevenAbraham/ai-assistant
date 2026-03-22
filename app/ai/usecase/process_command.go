package usecase

import (
	"context"
	"fmt"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/app/ai/repository"
)

// AIClient is the interface for calling an LLM backend (e.g. Claude).
// Implemented in internal/httpclient/.
type AIClient interface {
	Complete(ctx context.Context, messages []entity.Message) (string, error)
}

// ProcessCommandInput carries the raw text and session context.
type ProcessCommandInput struct {
	Text      string
	SessionID string
}

// ProcessCommandOutput holds the assistant's response and detected intent.
type ProcessCommandOutput struct {
	Response string
	Intent   entity.Intent
}

// ProcessCommandUseCase orchestrates the main request-response loop.
type ProcessCommandUseCase struct {
	conversationRepo repository.ConversationRepository
	memoryRepo       repository.MemoryRepository
	aiClient         AIClient
}

func NewProcessCommandUseCase(
	convRepo repository.ConversationRepository,
	memRepo repository.MemoryRepository,
	aiClient AIClient,
) *ProcessCommandUseCase {
	return &ProcessCommandUseCase{
		conversationRepo: convRepo,
		memoryRepo:       memRepo,
		aiClient:         aiClient,
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

	responseText, err := uc.aiClient.Complete(ctx, conv.Messages)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", entity.ErrAIClientFailure, err)
	}

	assistantMsg := entity.Message{
		Role:    entity.RoleAssistant,
		Content: responseText,
	}
	conv.Messages = append(conv.Messages, assistantMsg)

	if saveErr := uc.conversationRepo.Save(ctx, conv); saveErr != nil {
		// Non-fatal: log but do not block the response.
		_ = saveErr
	}

	return &ProcessCommandOutput{
		Response: responseText,
		Intent:   entity.IntentChat,
	}, nil
}
