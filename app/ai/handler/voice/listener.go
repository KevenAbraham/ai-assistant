package voice

import (
	"context"
	"fmt"
	"math"

	"github.com/gordonklaus/portaudio"
)

const (
	sampleRate      = 16000
	framesPerBuffer = 1024
)

// AudioCapture is the interface for capturing raw audio samples.
// The concrete implementation lives in Listener.
type AudioCapture interface {
	Listen(ctx context.Context) ([]int16, error)
}

// ListenerConfig holds tunable parameters for the Listener.
type ListenerConfig struct {
	// MaxRecordSeconds is the hard upper limit on recording duration.
	MaxRecordSeconds int
	// SilenceThreshold is the RMS amplitude below which audio is considered silence (0–32767 scale).
	SilenceThreshold float64
	// SilenceDurationMs is how long continuous silence must last (in ms) before recording stops.
	SilenceDurationMs int
}

// Listener captures audio from the default input device using PortAudio.
// It stops recording automatically once the user stops speaking (VAD), or when
// MaxRecordSeconds is reached — whichever comes first.
type Listener struct {
	cfg ListenerConfig
}

func NewListener(cfg ListenerConfig) *Listener {
	return &Listener{cfg: cfg}
}

func (l *Listener) Listen(ctx context.Context) ([]int16, error) {
	if err := portaudio.Initialize(); err != nil {
		return nil, fmt.Errorf("portaudio initialize: %w", err)
	}
	defer portaudio.Terminate() //nolint:errcheck

	chunk := make([]int16, framesPerBuffer)
	stream, err := portaudio.OpenDefaultStream(1, 0, float64(sampleRate), framesPerBuffer, chunk)
	if err != nil {
		return nil, fmt.Errorf("open stream: %w", err)
	}
	defer stream.Close() //nolint:errcheck

	if err := stream.Start(); err != nil {
		return nil, fmt.Errorf("start stream: %w", err)
	}
	defer stream.Stop() //nolint:errcheck

	maxChunks := l.cfg.MaxRecordSeconds * sampleRate / framesPerBuffer
	silenceChunkThreshold := l.cfg.SilenceDurationMs * sampleRate / (1000 * framesPerBuffer)

	var all []int16
	silentChunks := 0
	hasSpeech := false

	for i := 0; i < maxChunks; i++ {
		select {
		case <-ctx.Done():
			return all, nil
		default:
		}

		if err := stream.Read(); err != nil {
			return nil, fmt.Errorf("read stream: %w", err)
		}

		copied := make([]int16, framesPerBuffer)
		copy(copied, chunk)
		all = append(all, copied...)

		if rmsAmplitude(chunk) >= l.cfg.SilenceThreshold {
			hasSpeech = true
			silentChunks = 0
			continue
		}

		if hasSpeech {
			silentChunks++
			if silentChunks >= silenceChunkThreshold {
				break
			}
		}
	}

	return all, nil
}

func rmsAmplitude(samples []int16) float64 {
	if len(samples) == 0 {
		return 0
	}
	var sum float64
	for _, s := range samples {
		sum += float64(s) * float64(s)
	}
	return math.Sqrt(sum / float64(len(samples)))
}
