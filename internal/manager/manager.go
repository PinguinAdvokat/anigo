package manager

import (
	"anigo/internal/extractors"
	"log"
	"sync"
)

type Extractor interface {
	Search(string) ([]extractors.Anime, error)
	ParseAnime(extractors.Anime) (extractors.Anime, error)
	ParseEpisode(*extractors.Episode, string, string) error
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

func (m *Manager) ParseAnime(animeIndex int) error {
	parsedAnime, err := m.Extractor.ParseAnime(m.FoundAnime[animeIndex])
	if err != nil {
		log.Printf("manager failed in parsing anime: %v\n", err)
		return err
	}
	parsedAnime.Parsed = true
	m.FoundAnime[animeIndex] = parsedAnime
	return nil
}

func (m *Manager) ParseEpisode(animeIndex, episodeIndex int, player, voicecover string) error {
	return m.Extractor.ParseEpisode(&m.FoundAnime[animeIndex].Episodes[episodeIndex], player, voicecover)
}
