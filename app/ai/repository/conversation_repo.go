package repository

import (
	"context"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// ConversationRepository defines persistence operations for conversations.
// Implementations live in internal/repository/.
type ConversationRepository interface {
	Save(ctx context.Context, conv *entity.Conversation) error
	FindByID(ctx context.Context, id string) (*entity.Conversation, error)
	FindBySessionID(ctx context.Context, sessionID string) (*entity.Conversation, error)
	FindRecent(ctx context.Context, limit int) ([]*entity.Conversation, error)
	AppendMessage(ctx context.Context, conversationID string, msg entity.Message) error
}
