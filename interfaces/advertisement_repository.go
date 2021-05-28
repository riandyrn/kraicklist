package interfaces

import (
	"fmt"
	"github.com/knightazura/domain"
)

type AdvertisementRepository struct {}

func(ar *AdvertisementRepository) Store() (newAd domain.Advertisement, newDoc domain.GeneralDocument) {
	return
}

func (ar *AdvertisementRepository) BulkStore(ads []domain.Advertisement) (newAds []domain.Advertisement, newDocs domain.GeneralDocuments) {
	if len(ads) == 0 {
		_ = fmt.Errorf("There's no ads that need to be converted to docs")
		return
	}

	for _, ad := range ads {
		newAds = append(newAds, ad)
		newDocs = append(newDocs, domain.GeneralDocument{
			ID:      ad.ID,
			Title:   ad.Title,
			Content: ad.Content,
		})
	}
	return
}