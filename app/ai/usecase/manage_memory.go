package usecase

import (
	"context"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/app/ai/repository"
)

// ManageMemoryUseCase handles saving and retrieving long-term memories.
type ManageMemoryUseCase struct {
	memoryRepo repository.MemoryRepository
}

func NewManageMemoryUseCase(repo repository.MemoryRepository) *ManageMemoryUseCase {
	return &ManageMemoryUseCase{memoryRepo: repo}
}

func (uc *ManageMemoryUseCase) Save(ctx context.Context, key, value string) error {
	mem := &entity.Memory{
		Key:   key,
		Value: value,
	}
	return uc.memoryRepo.Save(ctx, mem)
}

func (uc *ManageMemoryUseCase) Search(ctx context.Context, query string) ([]*entity.Memory, error) {
	return uc.memoryRepo.Search(ctx, query)
}

func (uc *ManageMemoryUseCase) FindAll(ctx context.Context) ([]*entity.Memory, error) {
	return uc.memoryRepo.FindAll(ctx)
}

func (uc *ManageMemoryUseCase) Delete(ctx context.Context, key string) error {
	return uc.memoryRepo.Delete(ctx, key)
}
