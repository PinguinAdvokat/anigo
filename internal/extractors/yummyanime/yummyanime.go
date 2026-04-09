package yummyanime

import (
	"anigo/internal/parsers/kodik"
	"net/http"
)

type YummyAnime struct {
	BaseURL     string
	httpClient  *http.Client
	kodikParser *kodik.Kodik
}

func NewYummyAnime(kodikParser *kodik.Kodik, httpClient *http.Client) *YummyAnime {
	return &YummyAnime{
		BaseURL:     "https://api.yani.tv",
		httpClient:  httpClient,
		kodikParser: kodikParser,
	}
}
