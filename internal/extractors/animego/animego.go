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

func New(httpClient *http.Client) *Animego {
	return &Animego{
		BaseURL:    "https://animego.me",
		httpClient: httpClient,
	}
}
