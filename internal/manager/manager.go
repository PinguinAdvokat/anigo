package manager

import (
	"anigo/internal/extractors"
	"log"
	"sync"
)

type Extractor interface {
	Search(string) ([]extractors.Anime, error)
}

type Manager struct {
	mu         sync.RWMutex
	Extractor  Extractor
	FoundAnime []extractors.Anime
}

func New(extractor Extractor) *Manager {
	return &Manager{Extractor: extractor}
}

func (m *Manager) Search(name string) error {
	log.Println("searching " + name)
	animes, err := m.Extractor.Search(name)
	if err != nil {
		log.Printf("error in search: %s\n", err)
		return err
	}
	m.mu.Lock()
	m.FoundAnime = animes
	m.mu.Unlock()
	return err
}
