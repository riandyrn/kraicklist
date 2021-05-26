package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/knightazura/data/model"
	"github.com/knightazura/data/repository"
	"github.com/knightazura/vendors"
	"net/http"
)

// A wrapper service of Search Engine
// regardless of which search engine was used by the app
type Engine struct {
	LocalDocuments 			model.GeneralDocuments
	MeilisearchDocumets 	model.MeilisearchDocuments
}

const activeEngine = "local"

// Handler of HTTP request for search endpoint
func (e *Engine) HandleSearch() http.HandlerFunc {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			// fetch query string from query params
			q := r.URL.Query().Get("q")
			if len(q) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte("missing search query in query params"))
				return
			}
			// search relevant records
			records, err := e.Search(q)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}
			// output success response
			buf := new(bytes.Buffer)
			encoder := json.NewEncoder(buf)
			encoder.Encode(records)
			w.Header().Set("Content-Type", "application/json")
			w.Write(buf.Bytes())
		},
	)
}

// Abstraction method of "search"
func (e *Engine) Search(query string) (docs []model.SearchResponse, err error) {
	switch activeEngine {
	case "local":
		docs = vendors.LocalPerformSearch(query, &e.LocalDocuments)
	case "default":
		_ = fmt.Errorf("No search engine vendor")
	}
	return
}

// Abstraction method of indexing model to be document
func (e *Engine) SetupDocument(docs model.GeneralDocuments, docType string, engineType string) {
	i := repository.IndexedDocument{}

	switch engineType {
	case "meilisearch":
		e.MeilisearchDocumets = i.ToMeilisearchDocument(docs, docType)
	default:
		e.LocalDocuments = i.ToLocalDocument(docs)
	}
	return
}