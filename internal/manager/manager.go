package manager

import (
	"anigo/internal/extractors"
	"anigo/internal/parsers/kodik"
	"log"
	"sync"
)

type Extractor interface {
	Search(string) ([]extractors.Anime, error)
	ParseAnime(*extractors.Anime) error
	ParseEpisode(*extractors.Episode, string, string) error
}

type Manager struct {
	mu          sync.RWMutex
	Extractor   Extractor
	KodikParser *kodik.Kodik
	FoundAnime  []extractors.Anime
}

func New(extractor Extractor, kodikParser *kodik.Kodik) *Manager {
	return &Manager{Extractor: extractor, KodikParser: kodikParser}
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
	anime := &m.FoundAnime[animeIndex]
	err := m.Extractor.ParseAnime(anime)
	if err != nil {
		log.Printf("manager failed in parsing anime: %v\n", err)
		return err
	}
	anime.Parsed = true
	return nil
}

func (m *Manager) ParseEpisode(animeIndex, episodeIndex int, player, voicecover string) error {
	episode := &m.FoundAnime[animeIndex].Episodes[episodeIndex]
	err := m.Extractor.ParseEpisode(episode, player, voicecover)
	if err != nil {
		return err
	}

	links, err := m.KodikParser.Parse(episode.PlayerURL)
	if err != nil {
		return err
	}
	episode.Links = links
	return nil
}
