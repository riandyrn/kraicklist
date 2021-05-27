package services

// import (
// 	"bytes"
// 	"encoding/json"
// 	"fmt"
// 	"github.com/knightazura/contracts"
// 	"github.com/knightazura/data/model"
// 	"github.com/knightazura/data/repository"
// 	"github.com/knightazura/vendors"
// 	"log"
// 	"net/http"
// 	"os"
// )

// // A wrapper service of Search Engine
// // regardless of which search engine was used by the app
// type Engine struct {
// 	Default                 contracts.SearchEngine
// 	LocalDocuments 			model.GeneralDocuments
// }

// func InitSearchEngine(indexName string) (*Engine, error) {
// 	engine := getEngine(indexName)
// 	return &Engine{
// 		Default: engine,
// 	}, nil
// }

// // Handler of HTTP request for search endpoint
// func (e *Engine) HandleSearch() http.HandlerFunc {
// 	return http.HandlerFunc(
// 		func(w http.ResponseWriter, r *http.Request) {
// 			// fetch query string from query params
// 			q := r.URL.Query().Get("q")
// 			muid := r.URL.Query().Get("muid")
// 			if muid == "" {
// 				muid = "advertisement"
// 			}

// 			if len(q) == 0 {
// 				w.WriteHeader(http.StatusBadRequest)
// 				w.Write([]byte("missing search query in query params"))
// 				return
// 			}
// 			// search relevant records
// 			records, err := e.Search(q, muid)
// 			if err != nil {
// 				w.WriteHeader(http.StatusInternalServerError)
// 				w.Write([]byte(err.Error()))
// 				return
// 			}
// 			// output success response
// 			buf := new(bytes.Buffer)
// 			encoder := json.NewEncoder(buf)
// 			encoder.Encode(records)
// 			w.Header().Set("Content-Type", "application/json")
// 			w.Write(buf.Bytes())
// 		},
// 	)
// }

// // Abstraction method of "search"
// func (e *Engine) Search(query string, indexName string) (docs []model.SearchResponse, err error) {
// 	if activeEngine == "local" {
// 		docs = vendors.LocalPerformSearch(query, &e.LocalDocuments)
// 	}

// 	if activeEngine == "meilisearch" {
// 		client, err := vendors.InitMeilisearchEngine(indexName)
// 		if err != nil {
// 			fmt.Errorf(err.Error())
// 		}

// 		docs = client.PerformSearch(query)
// 	}
// 	_ = fmt.Errorf("No search engine vendor")
// 	return
// }

// // Abstraction method of indexing model to be document
// func (e *Engine) SetupDocument(docs model.GeneralDocuments, docType string, engineType string) {
// 	i := repository.IndexedDocument{}

// 	switch engineType {
// 	case "meilisearch":
// 		i.ToMeilisearchDocument(docs, docType)
// 	default:
// 		e.LocalDocuments = i.ToLocalDocument(docs)
// 	}
// 	return
// }

// // Get active search engine
// func getEngine(indexName string) contracts.SearchEngine {
// 	switch os.Getenv("SEARCH_ENGINE_ACTIVE") {
// 	case "meilisearch":
// 		engine, err := vendors.InitMeilisearchEngine(indexName)
// 		if err != nil {
// 			log.Fatalln("Failed to initialize Meilisearch instance")
// 		}
// 		return engine
// 	}
// 	return nil
// }