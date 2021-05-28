package contracts

import (
	"github.com/knightazura/domain"
)

type SearchEngine interface {
	PerformSearch(query string, indexName string) []domain.IndexedDocument
	IndexDocuments(docs domain.GeneralDocuments, indexName string)
}
