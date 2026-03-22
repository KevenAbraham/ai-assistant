package voice

import (
	"context"
)

// AudioTranscriber converts raw audio bytes to text.
// The concrete implementation lives in internal/httpclient/whisper_local_client.go.
type AudioTranscriber interface {
	Transcribe(ctx context.Context, audio []byte) (string, error)
}

// Transcriber is a voice-layer adapter that delegates to an AudioTranscriber.
type Transcriber struct {
	client AudioTranscriber
}

func NewTranscriber(client AudioTranscriber) *Transcriber {
	return &Transcriber{client: client}
}

func (t *Transcriber) Transcribe(ctx context.Context, samples []int16) (string, error) {
	// Convert int16 PCM to raw bytes (little-endian).
	raw := make([]byte, len(samples)*2)
	for i, s := range samples {
		raw[i*2] = byte(s)
		raw[i*2+1] = byte(s >> 8)
	}
	return t.client.Transcribe(ctx, raw)
}
