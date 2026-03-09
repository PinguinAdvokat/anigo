package extractors

type Anime struct {
	Title         string
	Rating        string
	URL           string
	EpisodesCount int
	Episodes      []Episode
}

type Episode struct {
	Name     string
	KodikURL string
}
