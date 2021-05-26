package vendors

import (
	"github.com/knightazura/data/model"
	"strings"
)

func LocalPerformSearch(query string, docs *model.GeneralDocuments) (result []model.SearchResponse) {
	for _, doc := range *docs {
		if strings.Contains(doc.Title, query) || strings.Contains(doc.Content, query) {
			result = append(result, model.SearchResponse{
				ID: doc.ID,
				Data: doc,
			})
		}
	}
	return
}
