package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/KevenAbraham/ai-assistant/app/ai/handler/voice"
	"github.com/KevenAbraham/ai-assistant/app/ai/service"
	"github.com/KevenAbraham/ai-assistant/app/ai/usecase"
	"github.com/KevenAbraham/ai-assistant/internal/config"
	"github.com/KevenAbraham/ai-assistant/internal/database"
	"github.com/KevenAbraham/ai-assistant/internal/httpclient"
	"github.com/KevenAbraham/ai-assistant/internal/repository"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	db, err := database.NewDB(ctx, cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer db.Close(ctx) //nolint:errcheck

	convRepo := repository.NewConversationRepository(db)
	memRepo := repository.NewMemoryRepository(db)
	var claudeClient usecase.AIClient = httpclient.NewClaudeClient(cfg)
	var whisperClient voice.AudioTranscriber = httpclient.NewWhisperClient(cfg)
	var ttsClient voice.TextSynthesizer = httpclient.NewTTSClient(cfg)

	systemPrompt, err := config.LoadSystemPrompt(cfg)
	if err != nil {
		log.Fatalf("system prompt: %v", err)
	}
	contextBuilder := service.NewContextBuilder(systemPrompt)
	actionExecutor := service.NewActionExecutor()

	var processCmd usecase.CommandProcessor = usecase.NewProcessCommandUseCase(convRepo, memRepo, claudeClient, contextBuilder, actionExecutor)

	var listener voice.AudioCapture = voice.NewListener(voice.ListenerConfig{
		MaxRecordSeconds:  cfg.RecordSeconds,
		SilenceThreshold:  cfg.SilenceThreshold,
		SilenceDurationMs: cfg.SilenceDurationMs,
	})
	transcriber := voice.NewTranscriber(whisperClient)
	speaker := voice.NewSpeaker(ttsClient)

	log.Println("voice daemon started — listening...")

	for {
		select {
		case <-ctx.Done():
			log.Println("shutting down voice daemon")
			return
		default:
		}

		samples, err := listener.Listen(ctx)
		if err != nil {
			log.Printf("listen error: %v", err)
			continue
		}

		text, err := transcriber.Transcribe(ctx, samples)
		if err != nil {
			log.Printf("transcribe error: %v", err)
			continue
		}
		if text == "" {
			continue
		}
		log.Printf("user: %s", text)

		out, err := processCmd.Execute(ctx, usecase.ProcessCommandInput{
			Text:      text,
			SessionID: "daemon",
		})
		if err != nil {
			log.Printf("process command error: %v", err)
			continue
		}
		log.Printf("assistant: %s", out.Response)
		
		if err := speaker.Speak(ctx, out.Response); err != nil {
			log.Printf("speak error: %v", err)
		}
	}
}
