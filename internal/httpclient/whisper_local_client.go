package httpclient

import (
	"bytes"
	"context"
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

func (c *WhisperLocalClient) Transcribe(ctx context.Context, audio []byte) (string, error) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)

	part, err := w.CreateFormFile("file", "audio.wav")
	if err != nil {
		return "", fmt.Errorf("whisper: create form file: %w", err)
	}
	if _, err := part.Write(audio); err != nil {
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
