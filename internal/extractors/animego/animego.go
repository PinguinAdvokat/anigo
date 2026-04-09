package animego

import (
	"net/http"
)

type Animego struct {
	BaseURL    string
	httpClient *http.Client
}

func New(httpClient *http.Client) *Animego {
	return &Animego{
		BaseURL:    "https://animego.me",
		httpClient: httpClient,
	}
}
