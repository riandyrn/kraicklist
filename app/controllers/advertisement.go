package controllers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/knightazura/contracts"
)

type Advertisement struct {
	EntityName string
	SearchEngine contracts.SearchEngine
}

const entityName = "advertisement"

// TO DO: to pass config value as parameter
func InitAdvertisement(se contracts.SearchEngine) *Advertisement {
	return &Advertisement{
		EntityName: entityName,
		SearchEngine: se,
	}
}

func (controller *Advertisement) Search(writer http.ResponseWriter, req *http.Request) {
	//context := req.Context()

	// Process the request
	query := req.URL.Query().Get("q")

	if len(query) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("missing search query in query params"))
		return
	}

	// Search relevant records
	searchEngine := controller.SearchEngine

	records := searchEngine.PerformSearch(query)
	// if err != nil {
	// 	writer.WriteHeader(http.StatusInternalServerError)
	// 	writer.Write([]byte(err.Error()))
	// 	return
	// }

	// output success response
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(records)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(buf.Bytes())
}