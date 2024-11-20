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
	// Initialize AWS client
	client, err := utils.NewEC2Client()
	if err != nil {
		log.Fatalf("failed to create EC2 client: %v", err)
	}

	// Initialize services and handlers
	ec2Service := &services.EC2Service{Client: client}
	ec2Handler := &handlers.EC2Handler{Service: ec2Service}

	// Initialize the router
	r := router.NewRouter(ec2Handler)

	// Start the server
	log.Println("Server is running on port 8080...")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	}
}
