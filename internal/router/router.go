package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/turaneminli/go_backend_aws/internal/handlers"
)

// NewRouter initializes and returns a new router
func NewRouter(ec2Handler *handlers.EC2Handler) http.Handler {
	r := chi.NewRouter()

	r.Get("/regions", ec2Handler.ListRegionsHandler)

	// instances
	r.Post("/instances/launch", ec2Handler.LaunchInstanceHandler)
	r.Post("/instances/stop", ec2Handler.StopInstanceByIdHandler)

	return r
}
