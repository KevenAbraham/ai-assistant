package usecase

import (
	"context"

	"github.com/KevenAbraham/ai-assistant/app/ai/entity"
	"github.com/KevenAbraham/ai-assistant/app/ai/repository"
)

// ManageHistoryUseCase handles reading conversation history.
type ManageHistoryUseCase struct {
	conversationRepo repository.ConversationRepository
}

func NewManageHistoryUseCase(repo repository.ConversationRepository) *ManageHistoryUseCase {
	return &ManageHistoryUseCase{conversationRepo: repo}
}

func (uc *ManageHistoryUseCase) GetBySession(ctx context.Context, sessionID string) (*entity.Conversation, error) {
	return uc.conversationRepo.FindBySessionID(ctx, sessionID)
}

func (uc *ManageHistoryUseCase) GetRecent(ctx context.Context, limit int) ([]*entity.Conversation, error) {
	return uc.conversationRepo.FindRecent(ctx, limit)
}
