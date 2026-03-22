package handler

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"github.com/mintrage/sysmon/internal/models"
	"github.com/mintrage/sysmon/internal/storage"
)

type Handler struct {
	Storage  *storage.Storage
	LastSeen map[string]time.Time
	Mu       sync.RWMutex
}

func (h *Handler) MetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST", http.StatusMethodNotAllowed)
		return
	}

	var m models.Metrics
	err := json.NewDecoder(r.Body).Decode(&m)
	if err != nil {
		http.Error(w, "Неверный формат JSON", http.StatusBadRequest)
		return
	}

	h.Mu.Lock()
	h.LastSeen[m.ServerName] = time.Now()
	h.Mu.Unlock()

	err = h.Storage.SaveMetric(m)
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LatestMetricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Только GET", http.StatusMethodNotAllowed)
		return
	}
	m, err := h.Storage.GetLatestMetric()
	if err != nil {
		http.Error(w, "500", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(m)
}
