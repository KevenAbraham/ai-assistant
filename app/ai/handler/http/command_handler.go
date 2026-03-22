package http

import (
	"encoding/json"
	"net/http"

	"github.com/KevenAbraham/ai-assistant/app/ai/usecase"
)

// CommandHandler handles POST /ai/command requests.
type CommandHandler struct {
	processCmd usecase.CommandProcessor
}

func NewCommandHandler(processCmd usecase.CommandProcessor) *CommandHandler {
	return &CommandHandler{processCmd: processCmd}
}

type commandRequest struct {
	Text      string `json:"text"`
	SessionID string `json:"session_id"`
}

type commandResponse struct {
	Response string `json:"response"`
	Intent   string `json:"intent"`
}

func (h *CommandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req commandRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	out, err := h.processCmd.Execute(r.Context(), usecase.ProcessCommandInput{
		Text:      req.Text,
		SessionID: req.SessionID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(commandResponse{ //nolint:errcheck
		Response: out.Response,
		Intent:   string(out.Intent),
	})
}
