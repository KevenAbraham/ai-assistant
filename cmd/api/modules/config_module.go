package modules

import (
	"go.uber.org/fx"

	"github.com/KevenAbraham/ai-assistant/internal/config"
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.Load),
)
