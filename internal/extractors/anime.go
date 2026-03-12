package extractors

type Anime struct {
	Title              string
	Description        string
	Rating             string
	URL                string
	AvailableVoiceover []string
	AvailablePlayers   []string
	SelectedVoiceover  []string
	SelectedSource     []string
	EpisodesCount      int
	Episodes           []Episode
}

type Episode struct {
	Name      string
	SourceURL string
}

var (
	AvailablePlayers = []string{"Kodik"}
)
