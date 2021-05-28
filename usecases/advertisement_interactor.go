package usecases

import (
	"github.com/knightazura/domain"
)

type AdvertisementInteractor struct {
	AdvertisementRepository AdvertisementRepository
	IndexedDocumentRepository IndexedDocumentRepository
}

const EntityName = "advertisement"

//func (adInteractor *AdvertisementInteractor) Store()

func (adInteractor *AdvertisementInteractor) Search(query string) []domain.IndexedDocument {
	docs := adInteractor.IndexedDocumentRepository.SearchDocs(query, EntityName)

	return docs
}

func (adInteractor *AdvertisementInteractor) Upload(ads []domain.Advertisement) (newAds []domain.Advertisement, docs domain.GeneralDocuments) {
	newAds, docs = adInteractor.AdvertisementRepository.BulkStore(ads)
	return
}

// Convert advertisement data to search engine document
// Should add context as first parameter
func (adInteractor *AdvertisementInteractor) ConvertToIndexedDocuments(docs domain.GeneralDocuments) {
	adInteractor.IndexedDocumentRepository.ToIndexedDocument(docs, EntityName)
	return
}