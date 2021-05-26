package model

import "fmt"

type Advertisement struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	ThumbURL  string   `json:"thumb_url"`
	Tags      []string `json:"tags"`
	UpdatedAt int64    `json:"updated_at"`
	ImageURLs []string `json:"image_urls"`
}

// Every model that want to be indexed,
// need to implement this method
func (ad *Advertisement) ConvertToGeneralDocs() (out GeneralDocument, err error) {
	if ad == nil {
		err = fmt.Errorf("There's no ads that need to be converted to docs")
		return
	}

	return GeneralDocument{
		ID: ad.ID,
		Title: ad.Title,
		Content: ad.Content,
		Model: "ads",
	}, nil
}