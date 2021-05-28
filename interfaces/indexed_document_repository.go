package interfaces

import (
	"github.com/knightazura/contracts"
	"github.com/knightazura/domain"
)

type IndexedDocumentRepository struct {
	SearchEngine contracts.SearchEngine
}

func (id *IndexedDocumentRepository) SearchDocs(query string, indexName string) []domain.IndexedDocument {
	// Deciding search engine vendor happened here
	docs := id.SearchEngine.PerformSearch(query, indexName)

	return docs
}

// Convert general document to meilisearch document
func (id *IndexedDocumentRepository) ToIndexedDocument(docs domain.GeneralDocuments, indexName string) {
	id.SearchEngine.IndexDocuments(docs, indexName)

	return
}