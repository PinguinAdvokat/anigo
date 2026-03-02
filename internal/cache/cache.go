package cache

import (
	"log"
	"sync"
)

type Cache struct {
	FilePath string
	sync.Map
}

func New(FilePath string) *Cache {
	cache := &Cache{FilePath: FilePath}
	if err := cache.LoadFromFile(); err != nil {
		log.Printf("cant load cache file: %v\n", err)
	}
	return cache
}
