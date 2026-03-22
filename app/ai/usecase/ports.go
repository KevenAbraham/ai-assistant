package usecase

import (
	"context"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// CommandProcessor is the input port for processing voice/text commands.
// Implemented by ProcessCommandUseCase.
type CommandProcessor interface {
	Execute(ctx context.Context, input ProcessCommandInput) (*ProcessCommandOutput, error)
}

// HistoryManager is the input port for reading conversation history.
// Implemented by ManageHistoryUseCase.
type HistoryManager interface {
	GetBySession(ctx context.Context, sessionID string) (*entity.Conversation, error)
	GetRecent(ctx context.Context, limit int) ([]*entity.Conversation, error)
}

// MemoryManager is the input port for managing long-term memories.
// Implemented by ManageMemoryUseCase.
type MemoryManager interface {
	Save(ctx context.Context, key, value string) error
	Search(ctx context.Context, query string) ([]*entity.Memory, error)
	FindAll(ctx context.Context) ([]*entity.Memory, error)
	Delete(ctx context.Context, key string) error
}
