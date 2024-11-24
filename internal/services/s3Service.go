package services

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// BucketInfo holds the information about each bucket
type BucketInfo struct {
	Name         string `json:"name"`
	Region       string `json:"region,omitempty"`
	CreationDate string `json:"creation_date"`
}

// S3Service is the service struct that holds the S3 client
type S3Service struct {
	Client *s3.Client
	cache  map[string]string // Cache for bucket regions
}

// NewS3Service initializes the S3Service with a cache
func NewS3Service(client *s3.Client) *S3Service {
	return &S3Service{
		Client: client,
		cache:  make(map[string]string),
	}
}

// ListBuckets retrieves a list of all buckets from S3 without fetching the regions by default
func (s *S3Service) ListBuckets() ([]BucketInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	output, err := s.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	var buckets []BucketInfo
	for _, bucket := range output.Buckets {
		buckets = append(buckets, BucketInfo{
			Name:         aws.ToString(bucket.Name),
			CreationDate: bucket.CreationDate.String(),
		})
	}

	return buckets, nil
}

// GetBucketRegion retrieves the region of the specified bucket
func (s *S3Service) GetBucketRegion(bucketName string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Check cache first
	region, found := s.cache[bucketName]
	if found {
		log.Printf("Cache hit for bucket: %s", bucketName)
		return region, nil
	}

	log.Printf("Cache miss for bucket: %s", bucketName)
	for retries := 0; retries < 3; retries++ {
		region, err := getBucketRegion(ctx, s.Client, bucketName)
		if err == nil {
			// Update cache
			s.cache[bucketName] = region
			return region, nil
		}
		log.Printf("Retrying to get region for bucket %s, attempt %d", bucketName, retries+1)
		time.Sleep(time.Duration(retries) * time.Second)
	}

	return "", fmt.Errorf("failed to get region for bucket %s", bucketName)
}

// getBucketRegion retrieves the region of the specified bucket
func getBucketRegion(ctx context.Context, svc *s3.Client, bucketName string) (string, error) {
	resp, err := svc.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return "", err
	}

	if resp.LocationConstraint == types.BucketLocationConstraint("") {
		return "us-east-1", nil
	}

	return string(resp.LocationConstraint), nil
}
