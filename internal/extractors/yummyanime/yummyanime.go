package yummyanime

import (
	"net/http"
)

type YummyAnime struct {
	BaseURL    string
	httpClient *http.Client
}

func New(httpClient *http.Client) *YummyAnime {
	return &YummyAnime{
		BaseURL:    "https://api.yani.tv",
		httpClient: httpClient,
	}
}
