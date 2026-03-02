package manager

import "anigo/internal/extractors"

type Extractor interface {
	Search(string) ([]extractors.AnimeInfo, error)
}

type Manager struct {
	Extractor Extractor
}

func New(extractor Extractor) *Manager {
	return &Manager{Extractor: extractor}
}

func (m *Manager) Search(name string) ([]extractors.AnimeInfo, error) {
	return m.Extractor.Search(name)
}
