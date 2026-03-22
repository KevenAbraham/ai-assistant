package modules

import (
	"context"
	"net/http"

	"go.uber.org/fx"
	"go.uber.org/zap"

	httphandler "github.com/KevenAbraham/ai-assistant/app/ai/handler/http"
	"github.com/KevenAbraham/ai-assistant/internal/config"
)

var ServerModule = fx.Module("server",
	fx.Provide(zap.NewProduction),
	fx.Invoke(startServer),
)

func startServer(
	lc fx.Lifecycle,
	cfg *config.Config,
	log *zap.Logger,
	cmdHandler *httphandler.CommandHandler,
	histHandler *httphandler.HistoryHandler,
) {
	mux := http.NewServeMux()
	mux.Handle("/ai/command", cmdHandler)
	mux.Handle("/ai/history", histHandler)
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`)) //nolint:errcheck
	})

	srv := &http.Server{
		Addr:    cfg.HTTPAddr,
		Handler: mux,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Info("starting HTTP server", zap.String("addr", cfg.HTTPAddr))
			go func() {
				if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Error("HTTP server error", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Info("shutting down HTTP server")
			return srv.Shutdown(ctx)
		},
	})
}
