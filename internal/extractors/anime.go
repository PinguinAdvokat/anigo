package extractors

type Anime struct {
	Title              string
	Description        string
	CoverURL           string
	Rating             string
	URL                string
	AvailableVoiceover []string
	AvailablePlayers   []string
	SelectedVoiceover  []string
	SelectedSource     []string
	EpisodesCount      int
	Episodes           []Episode
	Parsed             bool

	//[player][voiceover]episode
	YummyApisodesRaw map[string]map[string][]Episode
}

type Episode struct {
	Title     string
	URL       string
	PlayerURL string
	Links     map[string]string
}

var (
	AvailablePlayers = []string{"Kodik"}
)
