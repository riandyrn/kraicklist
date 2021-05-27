package controllers

import (
	"bytes"
	"encoding/json"
	"github.com/knightazura/services"
	"net/http"
)

type Advertisement struct {
	Services AdServices
}

type AdServices struct {
	SearchEngine    *services.Engine
	Seeder          *services.Seeder
}

func InitAdvertisement(services AdServices) *Advertisement {
	return &Advertisement{
		Services: services,
	}
}

func (controller *Advertisement) Search(writer http.ResponseWriter, req *http.Request) {
	//context := req.Context()

	// Process the request
	query := req.URL.Query().Get("q")
	muid := req.URL.Query().Get("muid")
	if muid == "" {
		muid = "advertisement"
	}

	if len(query) == 0 {
		writer.WriteHeader(http.StatusBadRequest)
		writer.Write([]byte("missing search query in query params"))
		return
	}

	// Search relevant records
	searchEngine := controller.Services.SearchEngine

	records, err := searchEngine.Search(query, muid)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		writer.Write([]byte(err.Error()))
		return
	}

	// output success response
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	encoder.Encode(records)
	writer.Header().Set("Content-Type", "application/json")
	writer.Write(buf.Bytes())
}