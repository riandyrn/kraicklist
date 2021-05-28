package infrastructure

import (
	"fmt"
	"github.com/knightazura/interfaces"
	"github.com/knightazura/services"
	"log"
	"net/http"
	"os"

	"github.com/knightazura/contracts"
)

type Services struct{
	Seeder *services.Seeder
	SearchEngine contracts.SearchEngine
}

func Dispatch(/** services should be parameters here */) {
	services := setupServices()
	setupServer(services)
}

func setupServices() *Services {
	// Search engine
	searchEngine, _ := InitSearchEngine()

	return &Services{
		Seeder: &services.Seeder{},
		SearchEngine: searchEngine,
	}
}

func setupServer(services *Services) {
	// Handle static files for frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	adController := interfaces.InitAdvertisementController(services.SearchEngine, services.Seeder)

	// Advertisement routes
	http.HandleFunc("/advertisement/search", adController.Search)
	// Challenge purpose: mock of /advertisement/upload endpoint
	adController.Upload()

	// Setup and start server
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3001"
	}
	fmt.Printf("Server is listening on %s...", port)
	
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("unable to start server due: %v", err)
	}
}