package manager

import (
	"anigo/internal/extractors"
	"anigo/internal/extractors/animego"
	"anigo/internal/extractors/yummyanime"
	"anigo/internal/mpv"
	"anigo/internal/parsers/kodik"
	"log"
	"net/http"
	"sync"
	"time"
)

type Extractor interface {
	Search(string) ([]extractors.Anime, error)
	ParseAnime(*extractors.Anime) error
	ParseEpisode(*extractors.Episode, string, string) error
}

type ExtractorFactory struct {
	httpClient *http.Client
}

func (e *ExtractorFactory) New(kind string) Extractor {
	switch kind {
	case "animego":
		return animego.New(e.httpClient)
	case "yummyanime":
		return yummyanime.New(e.httpClient)
	default:
		return animego.New(e.httpClient)
	}
}

type Manager struct {
	mu          sync.RWMutex
	factory     ExtractorFactory
	Extractor   Extractor
	KodikParser *kodik.Kodik
	FoundAnime  []extractors.Anime
	Mpv         *mpv.Mpv
}

func New(ExtractorKind string, kodikParser *kodik.Kodik, httpClient *http.Client, mpv *mpv.Mpv) *Manager {
	m := &Manager{
		factory:     ExtractorFactory{httpClient: httpClient},
		KodikParser: kodikParser,
		Mpv:         mpv,
	}
	m.SetExtractor(ExtractorKind)
	return m
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

func (m *Manager) PlayEpisode(animeIdx, episodeIdx int, quality string) {
	anime := m.FoundAnime[animeIdx]
	go func() {
		var url string
		for i := episodeIdx; i < len(anime.Episodes); i++ {
			m.ParseEpisode(animeIdx, i)
			url = anime.Episodes[i].Links[quality]
			log.Printf("url for playlist: %s", url)
			m.Mpv.Play(url)
			time.Sleep(time.Second * 3)
		}
	}()
}

func (m *Manager) SetExtractor(kind string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.FoundAnime = nil
	m.Extractor = m.factory.New(kind)
}
