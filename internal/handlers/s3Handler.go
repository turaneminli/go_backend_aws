package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/turaneminli/go_backend_aws/internal/services"
)

type S3Handler struct {
	Service *services.S3Service
}

type BucketResponse struct {
	Message string `json:"message"`
	Bucket  string `json:"bucket"`
}

func (h *S3Handler) ListBucketsHandler(w http.ResponseWriter, r *http.Request) {
	buckets, err := h.Service.ListBuckets()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(buckets); err != nil {
		http.Error(w, "Failed to encode buckets to JSON", http.StatusInternalServerError)
	}
}

func (h *S3Handler) CreateBucketHandler(w http.ResponseWriter, r *http.Request) {
	var input services.CreateBucketInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.Service.CreateBucket(input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := BucketResponse{
		Message: "Bucket created successfully",
		Bucket:  input.BucketName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func (h *S3Handler) DeleteBucketHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucketName")

	if bucketName == "" {
		http.Error(w, "Missing bucket name", http.StatusBadRequest)
		return
	}

	err := h.Service.DeleteBucket(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := BucketResponse{
		Message: "Bucket deleted successfully",
		Bucket:  bucketName,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *S3Handler) ListBucketObjectsHandler(w http.ResponseWriter, r *http.Request) {
	bucketName := r.URL.Query().Get("bucketName")

	if bucketName == "" {
		http.Error(w, "Missing bucket name", http.StatusBadRequest)
		return
	}

	objects, err := h.Service.ListBucketObjects(bucketName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(objects); err != nil {
		http.Error(w, "Failed to encode objects to JSON", http.StatusInternalServerError)
	}
}
