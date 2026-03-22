package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/KevenAbraham/ai-assistant/app/ai/usecase"
)

// HistoryHandler handles GET /ai/history requests.
type HistoryHandler struct {
	manageHistory usecase.HistoryManager
}

func NewHistoryHandler(manageHistory usecase.HistoryManager) *HistoryHandler {
	return &HistoryHandler{manageHistory: manageHistory}
}

// ServeHTTP handles GET /ai/history?session_id=xxx or ?limit=N.
func (h *HistoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	sessionID := r.URL.Query().Get("session_id")
	w.Header().Set("Content-Type", "application/json")

	if sessionID != "" {
		conv, err := h.manageHistory.GetBySession(r.Context(), sessionID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(conv) //nolint:errcheck
		return
	}

	limit := 10
	if l := r.URL.Query().Get("limit"); l != "" {
		if n, err := strconv.Atoi(l); err == nil && n > 0 {
			limit = n
		}
	}

	convs, err := h.manageHistory.GetRecent(r.Context(), limit)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(convs) //nolint:errcheck
}
