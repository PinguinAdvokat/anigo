package kodik

import (
	"anigo/internal/cache"
	"net/http"
)

type Kodik struct {
	HttpClient *http.Client
	cache      *cache.Cache
}

func New(httpClient *http.Client, cache *cache.Cache) *Kodik {
	return &Kodik{HttpClient: httpClient, cache: cache}
}
