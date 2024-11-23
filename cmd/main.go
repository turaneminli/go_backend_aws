package main

import (
	"log"
	"net/http"

	"github.com/turaneminli/go_backend_aws/internal/handlers"
	"github.com/turaneminli/go_backend_aws/internal/router"
	"github.com/turaneminli/go_backend_aws/internal/services"
	"github.com/turaneminli/go_backend_aws/internal/utils"
)

func main() {
	// Initialize EC2 client and service
	ec2Client, err := utils.NewEC2Client()
	if err != nil {
		log.Fatalf("failed to create EC2 client: %v", err)
	}
	ec2Service := &services.EC2Service{Client: ec2Client}
	ec2Handler := &handlers.EC2Handler{Service: ec2Service}

	// Initialize CloudWatch client and service
	cloudWatchClient, err := utils.CreateCloudWatchClient()
	if err != nil {
		log.Fatalf("failed to create CloudWatch client: %v", err)
	}
	cloudWatchService := &services.CloudWatchService{Client: cloudWatchClient}
	cloudWatchHandler := &handlers.CloudWatchHandler{Service: cloudWatchService}

	// Initialize S3 client and service
	s3Client, err := utils.NewS3Client()
	if err != nil {
		log.Fatalf("failed to create S3 client: %v", err)
	}
	s3Service := &services.S3Service{Client: s3Client}
	s3Handler := &handlers.S3Handler{Service: s3Service}

	// Initialize the router
	r := router.NewRouter(ec2Handler, cloudWatchHandler, s3Handler)

	log.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
