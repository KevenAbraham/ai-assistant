package repository

import (
	"context"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
)

// MemoryRepository defines persistence operations for long-term memories.
// Implementations live in internal/repository/.
type MemoryRepository interface {
	Save(ctx context.Context, mem *entity.Memory) error
	FindByKey(ctx context.Context, key string) (*entity.Memory, error)
	FindAll(ctx context.Context) ([]*entity.Memory, error)
	Search(ctx context.Context, query string) ([]*entity.Memory, error)
	Delete(ctx context.Context, key string) error
}
