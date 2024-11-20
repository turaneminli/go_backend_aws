package utils

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
)

func CreateCloudWatchClient() (*cloudwatch.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("unable to load AWS config: %v", err)
		return nil, err
	}

	// Create and return the CloudWatch client
	cloudWatchClient := cloudwatch.NewFromConfig(cfg)
	return cloudWatchClient, nil
}
