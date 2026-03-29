package usecase

import (
	"context"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

type CommandProcessor interface {
	Execute(ctx context.Context, input ProcessCommandInput) (*ProcessCommandOutput, error)
}

type HistoryManager interface {
	GetBySession(ctx context.Context, sessionID string) (*entity.Conversation, error)
	GetRecent(ctx context.Context, limit int) ([]*entity.Conversation, error)
}

type MemoryManager interface {
	Save(ctx context.Context, key, value string) error
	Search(ctx context.Context, query string) ([]*entity.Memory, error)
	FindAll(ctx context.Context) ([]*entity.Memory, error)
	Delete(ctx context.Context, key string) error
}
