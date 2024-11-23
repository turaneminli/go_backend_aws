package services

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3Service encapsulates S3 operations
type S3Service struct {
	Client *s3.Client
}

type BucketInfo struct {
	Name         string `json:"name"`
	CreationDate string `json:"creation_date"`
}

type CreateBucketInput struct {
	BucketName string `json:"bucket_name"`
	Region     string `json:"region"`
}

// ListBuckets retrieves a list of all buckets
func (s *S3Service) ListBuckets() ([]BucketInfo, error) {
	output, err := s.Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
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

// CreateBucket creates a new S3 bucket in the specified region
func (s *S3Service) CreateBucket(input CreateBucketInput) error {
	_, err := s.Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(input.BucketName),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraint(input.Region),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create bucket: %w", err)
	}

	return nil
}

// DeleteBucket deletes the specified S3 bucket
func (s *S3Service) DeleteBucket(bucketName string) error {
	_, err := s.Client.DeleteBucket(context.TODO(), &s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}

// GetBucketDetails retrieves the details of a specific bucket, including its region
func (s *S3Service) GetBucketDetails(bucketName string) (string, error) {
	output, err := s.Client.GetBucketLocation(context.TODO(), &s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return "", fmt.Errorf("failed to get bucket location: %w", err)
	}

	region := string(output.LocationConstraint)
	if region == "" {
		region = "us-east-1" // Default region if no location constraint is set
	}

	return region, nil
}

// ListBucketObjects lists objects in a specified bucket
func (s *S3Service) ListBucketObjects(bucketName string) ([]string, error) {
	output, err := s.Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list objects in bucket: %w", err)
	}

	var objectKeys []string
	for _, object := range output.Contents {
		objectKeys = append(objectKeys, aws.ToString(object.Key))
	}

	return objectKeys, nil
}
