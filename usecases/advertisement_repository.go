package usecases

import "github.com/knightazura/domain"

type AdvertisementRepository interface {
	Store() (newAd domain.Advertisement, newDoc domain.GeneralDocument)
	BulkStore(ads []domain.Advertisement) (newAds []domain.Advertisement, newDocs domain.GeneralDocuments)
}