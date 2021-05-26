package model

/** Constructor for any formatted documents by any vendors */

// Response format of searched documents
type SearchResponse struct {
	ID int64
	Data interface{}
}

// For local Search Engine
type GeneralDocument struct {
	ID        int64    `json:"id"`
	Title     string   `json:"title"`
	Content   string   `json:"content"`
	Model string `json:"model"`
}

type GeneralDocuments []GeneralDocument

// For vendor A Search Engine
type MeilisearchDocument struct {
	ID int64 `json:"id"`
	Data interface{}
}

type MeilisearchDocuments []MeilisearchDocument

// For vendor B Search Engine