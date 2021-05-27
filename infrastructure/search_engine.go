package infrastructure

import (
	"fmt"

	"github.com/knightazura/contracts"
	"github.com/knightazura/data/model"
	"github.com/knightazura/services"
	"github.com/knightazura/vendors"
)

type SearchEngineHandler struct {
	IndexName string
	Meilisearch *vendors.Meilisearch
	// another client can be added here
}

func InitSearchEngine(indexName string) (contracts.SearchEngine, error) {
	seHandler := &SearchEngineHandler{}

	// Meilisearch
	meilisearch := vendors.InitMeilisearch(indexName)

	seHandler.Meilisearch = meilisearch
	seHandler.IndexName = indexName

	return seHandler, nil
}

// Challenge purpose
func (se *SearchEngineHandler) SetupDocument() {
	seeder := &services.Seeder{}
	adDocs := seeder.LoadData("./data/data.gz", &model.Advertisement{})

	vendors.MSAddDocuments(se.Meilisearch.Client, adDocs, se.IndexName)
}

func (se *SearchEngineHandler) PerformSearch(query string) (docs []model.SearchResponse) {
	// TO DO: Check the config which search engine vendors is using
	docs = vendors.MSSearch(
		se.Meilisearch.Client, 
		se.IndexName, 
		query,
	)

	fmt.Print(se.IndexName, docs)
	return
}