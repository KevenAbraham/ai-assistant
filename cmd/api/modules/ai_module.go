package modules

import (
	"go.uber.org/fx"

	httphandler "github.com/KevenAbraham/ai-assistant/app/ai/handler/http"
	"github.com/KevenAbraham/ai-assistant/app/ai/service"
	"github.com/KevenAbraham/ai-assistant/app/ai/usecase"
	"github.com/KevenAbraham/ai-assistant/internal/config"
	"github.com/KevenAbraham/ai-assistant/internal/httpclient"
)

var AIModule = fx.Module("ai",
	// System prompt — reads file from disk and builds the ContextBuilder.
	fx.Provide(config.LoadSystemPrompt),
	fx.Provide(service.NewContextBuilder),

	// External API clients — ClaudeClient satisfies usecase.AIClient.
	fx.Provide(
		fx.Annotate(
			httpclient.NewClaudeClient,
			fx.As(new(usecase.AIClient)),
		),
	),

	// Use cases — registered as their input-port interfaces so handlers
	// depend on abstractions, not concrete types (Dependency Inversion).
	fx.Provide(
		fx.Annotate(
			usecase.NewProcessCommandUseCase,
			fx.As(new(usecase.CommandProcessor)),
		),
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewManageHistoryUseCase,
			fx.As(new(usecase.HistoryManager)),
		),
	),
	fx.Provide(
		fx.Annotate(
			usecase.NewManageMemoryUseCase,
			fx.As(new(usecase.MemoryManager)),
		),
	),

	// HTTP handlers.
	fx.Provide(httphandler.NewCommandHandler),
	fx.Provide(httphandler.NewHistoryHandler),
)
