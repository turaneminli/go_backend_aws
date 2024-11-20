package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/turaneminli/go_backend_aws/internal/services"
)

type EC2Handler struct {
	Service *services.EC2Service
}

type Response struct {
	Message    string `json:"message"`
	InstanceID string `json:"instance_id"`
}

// ListRegionsHandler handles the /regions endpoint
func (h *EC2Handler) ListRegionsHandler(w http.ResponseWriter, r *http.Request) {
	regions, err := h.Service.ListRegions()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"regions": regions,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// LaunchInstanceHandler handles the creation of new EC2 instances.
func (h *EC2Handler) LaunchInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var input services.LaunchInstanceInput

	// Parse the JSON request body
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Call the service layer to launch the instance
	instanceID, err := h.Service.LaunchInstance(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response object
	response := Response{
		Message:    "Instance launched successfully",
		InstanceID: instanceID,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Marshal the response into JSON
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *EC2Handler) StopInstanceByIdHandler(w http.ResponseWriter, r *http.Request) {
	instanceID := r.URL.Query().Get("instanceID")

	if instanceID == "" {
		http.Error(w, "Missing instance ID", http.StatusBadRequest)
		return
	}

	// Call the service layer to launch the instance
	instanceID, err := h.Service.StopInstanceById(instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Create the response object
	response := Response{
		Message:    "Instance stopped successfully",
		InstanceID: instanceID,
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Marshal the response into JSON
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}
