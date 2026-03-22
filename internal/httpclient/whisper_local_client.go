package httpclient

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/KevenAbraham/ai-assistant/internal/config"
)

// WhisperLocalClient sends audio to a local whisper.cpp server for transcription.
// Run the server with: whisper-server -m <model.bin> --port 9000
type WhisperLocalClient struct {
	serverURL  string
	httpClient *http.Client
}

func NewWhisperClient(cfg *config.Config) *WhisperLocalClient {
	return &WhisperLocalClient{
		serverURL:  cfg.WhisperURL,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// encodeWAV wraps raw 16-bit little-endian PCM (16 kHz, mono) in a WAV container.
func encodeWAV(pcm []byte) []byte {
	const (
		sampleRate    = 16000
		numChannels   = 1
		bitsPerSample = 16
	)
	byteRate := uint32(sampleRate * numChannels * bitsPerSample / 8)
	blockAlign := uint16(numChannels * bitsPerSample / 8)
	dataSize := uint32(len(pcm))

	var buf bytes.Buffer
	buf.WriteString("RIFF")
	binary.Write(&buf, binary.LittleEndian, 36+dataSize) //nolint:errcheck
	buf.WriteString("WAVE")
	buf.WriteString("fmt ")
	binary.Write(&buf, binary.LittleEndian, uint32(16))         //nolint:errcheck
	binary.Write(&buf, binary.LittleEndian, uint16(1))          //nolint:errcheck // PCM
	binary.Write(&buf, binary.LittleEndian, uint16(numChannels)) //nolint:errcheck
	binary.Write(&buf, binary.LittleEndian, uint32(sampleRate)) //nolint:errcheck
	binary.Write(&buf, binary.LittleEndian, byteRate)           //nolint:errcheck
	binary.Write(&buf, binary.LittleEndian, blockAlign)         //nolint:errcheck
	binary.Write(&buf, binary.LittleEndian, uint16(bitsPerSample)) //nolint:errcheck
	buf.WriteString("data")
	binary.Write(&buf, binary.LittleEndian, dataSize) //nolint:errcheck
	buf.Write(pcm)
	return buf.Bytes()
}

func (c *WhisperLocalClient) Transcribe(ctx context.Context, audio []byte) (string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", fmt.Errorf("whisper: create form file: %w", err)
	}
	if _, err := part.Write(encodeWAV(audio)); err != nil {
		return "", fmt.Errorf("whisper: write audio: %w", err)
	}
	w.Close()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.serverURL+"/inference", &buf)
	if err != nil {
		return "", fmt.Errorf("whisper: build request: %w", err)
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("whisper: request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("whisper: server returned %d: %s", resp.StatusCode, body)
	}

	var result struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("whisper: decode response: %w", err)
	}

	return result.Text, nil
}
