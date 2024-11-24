package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/cors"
	"github.com/turaneminli/go_backend_aws/internal/handlers"
)

// NewRouter initializes and returns a new router
func NewRouter(ec2Handler *handlers.EC2Handler, cloudWatchHandler *handlers.CloudWatchHandler, s3Handler *handlers.S3Handler) http.Handler {
	r := chi.NewRouter()

	// CORS configuration
	corsConfig := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                             // React app URL (adjust as needed)
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},        // Methods allowed
		AllowedHeaders:   []string{"Content-Type", "Authorization"}, // Allowed headers
		AllowCredentials: true,
		Debug:            true, // Enable debug to log CORS issues in the server logs
	})

	// Apply CORS middleware
	r.Use(corsConfig.Handler)

	// EC2 Routes
	r.Get("/regions", ec2Handler.ListRegionsHandler)
	r.Post("/instances/launch", ec2Handler.LaunchInstanceHandler)
	r.Post("/instances/stop", ec2Handler.StopInstanceByIdHandler)
	r.Post("/instances/start", ec2Handler.StartInstanceByIdHandler)
	r.Post("/instances/reboot", ec2Handler.RebootInstanceByIdHandler)
	r.Post("/instances/terminate", ec2Handler.TerminateInstanceByIdHandler)
	r.Get("/instances/status", ec2Handler.ListRunningInstancesStatusHandler)
	r.Get("/instances/detail", ec2Handler.InstanceDetailHandler)

	r.Get("/security-groups", ec2Handler.ListSecurityGroupsHandler)

	// CloudWatch Routes
	r.Get("/cloudwatch/metrics", cloudWatchHandler.GetEC2MetricsHandler)

	// S3 Routes
	r.Get("/s3/buckets", s3Handler.ListBucketsHandler)

	return r
}
