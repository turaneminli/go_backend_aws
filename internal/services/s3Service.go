package services

import (
	"context"
	"fmt"
	"sync"
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
}

// NewS3Service initializes the S3Service
func NewS3Service(client *s3.Client) *S3Service {
	return &S3Service{
		Client: client,
	}
}

// ListBuckets retrieves a list of all buckets and their regions
func (s *S3Service) ListBuckets() ([]BucketInfo, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Fetch the list of buckets
	output, err := s.Client.ListBuckets(ctx, &s3.ListBucketsInput{})
	if err != nil {
		return nil, fmt.Errorf("failed to list buckets: %w", err)
	}

	var wg sync.WaitGroup
	bucketsCh := make(chan BucketInfo, len(output.Buckets)) // Channel to collect results

	// Loop through each bucket and fetch the region concurrently using goroutines
	for _, bucket := range output.Buckets {
		wg.Add(1)
		go func(bucket types.Bucket) {
			defer wg.Done()
			region := s.getBucketRegion(ctx, aws.ToString(bucket.Name))
			bucketsCh <- BucketInfo{
				Name:         aws.ToString(bucket.Name),
				CreationDate: bucket.CreationDate.Format(time.RFC3339),
				Region:       region,
			}
		}(bucket)
	}

	// Wait for all goroutines to finish and close the channel
	go func() {
		wg.Wait()
		close(bucketsCh)
	}()

	var buckets []BucketInfo
	for bucket := range bucketsCh {
		buckets = append(buckets, bucket)
	}

	return buckets, nil
}

// getBucketRegion fetches the region for a bucket
func (s *S3Service) getBucketRegion(ctx context.Context, bucketName string) string {
	// Fetch the region from S3
	locationOutput, err := s.Client.GetBucketLocation(ctx, &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		fmt.Printf("failed to get location for bucket %s: %v\n", bucketName, err)
		return "Unknown"
	}

	// Default region if no location is found
	region := "us-east-1"
	if locationOutput.LocationConstraint != "" {
		region = string(locationOutput.LocationConstraint)
	}

	return region
}
