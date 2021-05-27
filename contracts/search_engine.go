package contracts

import "github.com/knightazura/data/model"

type SearchEngine interface {
	PerformSearch(query string) []model.SearchResponse
	SetupDocument()
}
