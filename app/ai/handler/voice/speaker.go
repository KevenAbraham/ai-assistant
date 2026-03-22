package voice

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// TextSynthesizer converts text to audio bytes.
// The concrete implementation lives in internal/httpclient/tts_local_client.go.
type TextSynthesizer interface {
	Synthesize(ctx context.Context, text string) ([]byte, error)
}

// Speaker converts text to speech and plays it through the default audio output.
type Speaker struct {
	tts TextSynthesizer
}

func NewSpeaker(tts TextSynthesizer) *Speaker {
	return &Speaker{tts: tts}
}

func (s *Speaker) Speak(ctx context.Context, text string) error {
	audio, err := s.tts.Synthesize(ctx, text)
	if err != nil {
		return fmt.Errorf("synthesis failed: %w", err)
	}

	f, err := os.CreateTemp("", "tts-*.wav")
	if err != nil {
		return fmt.Errorf("create temp file: %w", err)
	}
	defer os.Remove(f.Name()) //nolint:errcheck

	if _, err := f.Write(audio); err != nil {
		return fmt.Errorf("write audio: %w", err)
	}
	f.Close() //nolint:errcheck

	player := "aplay"
	if runtime.GOOS == "darwin" {
		player = "afplay"
	}
	return exec.CommandContext(ctx, player, f.Name()).Run()
}
