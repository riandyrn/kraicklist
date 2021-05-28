package vendors

import (
	"github.com/knightazura/domain"
	"log"
	"math/rand"

	"github.com/meilisearch/meilisearch-go"
)

type Meilisearch struct {
	Client meilisearch.ClientInterface
}

func InitMeilisearch() *Meilisearch {
	config := meilisearch.Config{
		Host: "http://127.0.0.1:7700",
	}

	client := meilisearch.NewClient(config)

	return &Meilisearch{
		Client: client,
	}
}

func MSAddDocuments(client meilisearch.ClientInterface, docs domain.GeneralDocuments, indexName string) {
	get, _ := client.Indexes().Get(indexName)

	// Create the index if it's not there
	if get == nil {
		_, err := client.Indexes().Create(meilisearch.CreateIndexRequest{
			UID: indexName,
		})

		if err != nil {
			log.Fatalf("Failed to create index of %s: %v", indexName, err)
			return
		}
	}

	var documents []domain.Advertisement
	for _, doc := range docs {
		documents = append(documents, domain.Advertisement{
			ID:      doc.ID,
			Title:   doc.Title,
			Content: doc.Content,
		})
	}

	_, err := client.Documents(indexName).AddOrUpdate(documents)
	if err != nil {
		log.Fatalf("Failed to add %s documents: %v", indexName, err)
		return
	}
	log.Println("Berhasil memasukkan dokumen")
}

func MSSearch(client meilisearch.ClientInterface, indexName string, query string) (docs []domain.IndexedDocument) {
	res, _ := client.Search(indexName).Search(meilisearch.SearchRequest{
		Query:  query,
		Limit:  10,
		Offset: 1,
	})

	for _, h := range res.Hits {
		docs = append(docs, domain.IndexedDocument{
			ID:   rand.Int63(),
			Data: h,
		})
	}
	return
}
