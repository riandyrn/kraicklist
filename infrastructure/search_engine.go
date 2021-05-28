package infrastructure

import (
	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
	"github.com/knightazura/vendors"
)

type SearchEngineHandler struct {
	IndexName string
	Meilisearch *vendors.Meilisearch
	// another client can be added here
}

func InitSearchEngine() (contracts.SearchEngine, error) {
	seHandler := &SearchEngineHandler{}

	// Create (any) search engine instances here
	seHandler.Meilisearch = vendors.InitMeilisearch()

	return seHandler, nil
}

func (se *SearchEngineHandler) IndexDocuments(docs domain.GeneralDocuments, indexName string) {
	// TO DO: Check the config which search engine client is using
	client := se.Meilisearch.Client
	vendors.MSAddDocuments(client, docs, indexName)
}

func (se *SearchEngineHandler) PerformSearch(query string, indexName string) (docs []domain.IndexedDocument) {
	// TO DO: Check the config which search engine client is using
	client := se.Meilisearch.Client

	docs = vendors.MSSearch(
		client,
		indexName,
		query,
	)
	return
}