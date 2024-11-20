package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/turaneminli/go_backend_aws/internal/services"
)

type CloudWatchHandler struct {
	Service *services.CloudWatchService
}

func (h *CloudWatchHandler) GetEC2MetricsHandler(w http.ResponseWriter, r *http.Request) {
	instanceID := r.URL.Query().Get("instanceId")
	if instanceID == "" {
		http.Error(w, "instance_id query parameter is required", http.StatusBadRequest)
		return
	}

	metrics, err := h.Service.GetEC2Metrics(instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		http.Error(w, "Failed to encode metrics to JSON", http.StatusInternalServerError)
	}
}
