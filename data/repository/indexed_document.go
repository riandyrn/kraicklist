package repository

import (
	"github.com/knightazura/data/model"
	"github.com/meilisearch/meilisearch-go"
	"log"
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
func (id *IndexedDocument) ToMeilisearchDocument(docs model.GeneralDocuments, indexName string) {
	var client = meilisearch.NewClient(meilisearch.Config{
		Host: "http://127.0.0.1:7700",
	})

	get, _ := client.Indexes().Get(indexName)

	// Create the index if it's not there
	if get == nil {
		_, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID: indexName,
		})

		if err != nil {
			log.Printf("Failed to create index of %s", indexName)
			return
		}
	}

	var documents []model.Advertisement
	for _, doc := range docs {
		documents = append(documents, model.Advertisement{
			ID: doc.ID,
			Title: doc.Title,
			Content: doc.Content,
		})
	}

	_, err := client.Documents(indexName).AddOrUpdate(documents)
	if err != nil {
		return
	}
	log.Println("Berhasil memasukkan dokumen")
	return
}