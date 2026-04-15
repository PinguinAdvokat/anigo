package kodik

import (
	"net/http"
)

type Kodik struct {
	HttpClient *http.Client
}

func New(httpClient *http.Client) *Kodik {
	return &Kodik{HttpClient: httpClient}
}
