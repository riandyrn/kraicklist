package repository

import (
	"github.com/knightazura/data/model"
)

type IndexedDocument struct {}

// Indexing for local vendor SE
func (id *IndexedDocument) ToLocalDocument(docs model.GeneralDocuments) (out model.GeneralDocuments){
	for _, doc := range docs {
		out = append(out, doc)
	}
	return
}

// Indexing for meilisearch vendor SE
func (id *IndexedDocument) ToMeilisearchDocument(docs model.GeneralDocuments, indexName string) (out model.MeilisearchDocuments) {
	_ = indexName
	for _, doc := range docs {
		ms := model.MeilisearchDocument{
			ID: doc.ID,
			Data: doc,
		}
		out = append(out, ms)
	}
	return
}