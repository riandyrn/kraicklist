package contracts

import (
	"github.com/knightazura/domain"
)

type SearchEngine interface {
	PerformSearch(query string) []model.SearchResponse
	SetupDocument()
}
