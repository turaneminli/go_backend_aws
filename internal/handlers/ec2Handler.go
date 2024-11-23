package handlers

import (
	"encoding/json"
	"fmt"
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

func (h *EC2Handler) LaunchInstanceHandler(w http.ResponseWriter, r *http.Request) {
	var input services.LaunchInstanceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	instanceID, err := h.Service.LaunchInstance(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message:    "Instance launched successfully",
		InstanceID: instanceID,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *EC2Handler) StopInstanceByIdHandler(w http.ResponseWriter, r *http.Request) {
	instanceID := r.URL.Query().Get("instanceId")

	if instanceID == "" {
		http.Error(w, "Missing instance ID", http.StatusBadRequest)
		return
	}

	instanceID, err := h.Service.StopInstanceById(instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message:    "Instance stopped successfully",
		InstanceID: instanceID,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *EC2Handler) StartInstanceByIdHandler(w http.ResponseWriter, r *http.Request) {
	instanceID := r.URL.Query().Get("instanceId")

	if instanceID == "" {
		http.Error(w, "Missing instance ID", http.StatusBadRequest)
		return
	}

	instanceID, err := h.Service.StartInstanceById(instanceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := Response{
		Message:    "Instance started successfully",
		InstanceID: instanceID,
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *EC2Handler) ListRunningInstancesStatusHandler(w http.ResponseWriter, r *http.Request) {
	instances, err := h.Service.GetAllRunningInstancesStatus()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(instances); err != nil {
		http.Error(w, "Failed to encode instances to JSON", http.StatusInternalServerError)
	}
}

func (h *EC2Handler) ListSecurityGroupsHandler(w http.ResponseWriter, r *http.Request) {
	securityGroups, err := h.Service.ListSecurityGroups()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to retrieve security groups: %v", err), http.StatusInternalServerError)
		return
	}

	// Return the list of security groups as JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(securityGroups); err != nil {
		http.Error(w, fmt.Sprintf("Failed to encode response: %v", err), http.StatusInternalServerError)
	}
}
