package modules

import (
	"context"

	"go.uber.org/fx"

	"github.com/KevenAbraham/ai-assistant/internal/config"
	"github.com/KevenAbraham/ai-assistant/internal/database"
	"github.com/KevenAbraham/ai-assistant/internal/repository"
)

// newDBFx is an fx-aware wrapper for database.NewDB.
func newDBFx(lc fx.Lifecycle, cfg *config.Config) (*database.DB, error) {
	db, err := database.NewDB(context.Background(), cfg)
	if err != nil {
		return nil, err
	}
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return db.Close(ctx)
		},
	})
	return db, nil
}

// DBModule provides the database connection and all repository implementations.
var DBModule = fx.Module("db",
	fx.Provide(newDBFx),
	fx.Provide(repository.NewConversationRepository),
	fx.Provide(repository.NewMemoryRepository),
)
