package vendors

import (
	"github.com/knightazura/data/model"
	"github.com/meilisearch/meilisearch-go"
	"math/rand"
)

type Meilisearch struct {
	Client      meilisearch.ClientInterface
	indexName   string
}

func InitMeilisearchEngine(indexName string) (*Meilisearch, error) {
	config := meilisearch.Config{
		Host: "http://127.0.0.1",
	}

	meilisearchInstance := meilisearch.NewClient(config)

	return &Meilisearch{
		Client: meilisearchInstance,
		indexName: indexName,
	}, nil
}

func (ms *Meilisearch) PerformSearch(query string) (docs []model.SearchResponse) {
	res, _ := ms.Client.Search(ms.indexName).Search(meilisearch.SearchRequest{
		Query: query,
		Limit: 10,
		Offset: 1,
	})

	for _, h := range res.Hits {
		docs = append(docs, model.SearchResponse{
			ID: rand.Int63(),
			Data: h,
		})
	}
	return
}
