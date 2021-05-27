package infrastructure

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/knightazura/app/controllers"
	"github.com/knightazura/contracts"
)

type Services struct{
	SearchEngine contracts.SearchEngine
}

func Dispatch(/** services should be parameters here */) {
	services := setupServices()
	setupServer(services)
}

func setupServices() *Services {
	// Search engine
	searchEngine, _ := InitSearchEngine("advertisement")
	
	// Setup document
	searchEngine.SetupDocument()

	return &Services{
		SearchEngine: searchEngine,
	}
}

func setupServer(services *Services) {
	// Handle static files for frontend
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	adController := controllers.InitAdvertisement(services.SearchEngine)

	// App routing and controllers
	http.HandleFunc("/search", adController.Search)

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