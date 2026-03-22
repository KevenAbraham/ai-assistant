package httpclient

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/KevenAbraham/ai-assistant/internal/config"
)

// LocalTTSClient synthesises speech using the OS built-in TTS engine.
// macOS: uses `say -o <file>` (built-in, no install needed)
// Linux: uses `espeak -w <file>` (install with: apt install espeak)
type LocalTTSClient struct{}

func NewTTSClient(_ *config.Config) *LocalTTSClient {
	return &LocalTTSClient{}
}

func (c *LocalTTSClient) Synthesize(ctx context.Context, text string) ([]byte, error) {
	ext := ".wav"
	if runtime.GOOS == "darwin" {
		ext = ".aiff"
	}

	f, err := os.CreateTemp("", "tts-*"+ext)
	if err != nil {
		return nil, fmt.Errorf("tts: create temp file: %w", err)
	}
	name := f.Name()
	f.Close()
	defer os.Remove(name)

	if out, err := ttsCommand(ctx, text, name).CombinedOutput(); err != nil {
		return nil, fmt.Errorf("tts: synthesis failed: %w — %s", err, out)
	}

	return os.ReadFile(name)
}

func ttsCommand(ctx context.Context, text, outputFile string) *exec.Cmd {
	if runtime.GOOS == "darwin" {
		return exec.CommandContext(ctx, "say", text, "-o", outputFile)
	}
	return exec.CommandContext(ctx, "espeak", text, "-w", outputFile)
}
