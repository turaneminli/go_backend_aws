package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/turaneminli/go_backend_aws/internal/services"
)

type S3Handler struct {
	Service *services.S3Service
}

// ListBucketsHandler handles the API request to get the list of buckets
func (h *S3Handler) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := h.Service.ListBuckets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")

	// Encode the list of buckets to JSON and send the response
	if err := json.NewEncoder(w).Encode(buckets); err != nil {
		http.Error(w, "Failed to encode buckets to JSON", http.StatusInternalServerError)
	}
}

// GetBucketRegionHandler handles the API request to get the region of a specific bucket
func (h *S3Handler) GetBucketRegionHandler(w http.ResponseWriter, r *http.Request) {
	// Extract bucket name from the URL path or query parameters
	bucketName := r.URL.Query().Get("bucketName")
	if bucketName == "" {
		http.Error(w, "Bucket name is required", http.StatusBadRequest)
		return
	}

	region, err := h.Service.GetBucketRegion(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set the response header to indicate JSON content
	w.Header().Set("Content-Type", "application/json")

	// Encode the region to JSON and send the response
	if err := json.NewEncoder(w).Encode(map[string]string{"region": region}); err != nil {
		http.Error(w, "Failed to encode region to JSON", http.StatusInternalServerError)
	}
}
