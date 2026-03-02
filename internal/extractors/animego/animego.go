package animego

import (
	"anigo/internal/parsers/kodik"
	"net/http"
)

type Animego struct {
	BaseURL     string
	httpClient  *http.Client
	kodikParser *kodik.Kodik
}

func New(kodikParser *kodik.Kodik, httpClient *http.Client) *Animego {
	return &Animego{
		BaseURL:     "https://animego.me",
		httpClient:  httpClient,
		kodikParser: kodikParser,
	}
}
