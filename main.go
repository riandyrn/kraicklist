package main

import (
	"fmt"
	"github.com/knightazura/app/controllers"
	"github.com/knightazura/data/model"
	"github.com/knightazura/services"
	"log"
	"net/http"
	"os"
)

func main() {
	// Setup controllers
	handler := controllers.Base{
		Advertisement: controllers.InitAdvertisement(),
	}
	// Init. Seeder
	seeder := services.Seeder{}

	// Seedr: Load Ads data to format as GeneralDocument
	ads := seeder.LoadData(
		"./data/data.gz",
		&model.Advertisement{},
	)

	// Initialize search engine
	searcher := &services.Engine{}

	// Setup Ads "GeneralDocument" to be indexed document
	searcher.SetupDocument(ads, "advertisement", "meilisearch")
	log.Printf("Jumlah data pada searcher %d dan ads %d", len(searcher.LocalDocuments), len(ads))

	// Define http controllers
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	http.HandleFunc("/search", searcher.HandleSearch())

	// define port, we need to set it as env for Heroku deployment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3001"
	}
	// start server
	fmt.Printf("Server is listening on %s...", port)
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
	if err != nil {
		log.Fatalf("unable to start server due: %v", err)
	}
}