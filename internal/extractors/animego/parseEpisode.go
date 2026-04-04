package animego

import (
	"anigo/internal/extractors"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"regexp"
)

type PlayerLinks map[string]map[string]string

var (
	reButton      = regexp.MustCompile(`(?i)<button\b([^>]+)>`)
	rePlayer      = regexp.MustCompile(`(?i)\bdata-player=["']([^"']+)["']`)
	reProvider    = regexp.MustCompile(`(?i)\bdata-provider-title=["']([^"']+)["']`)
	reTranslation = regexp.MustCompile(`(?i)\bdata-translation-title=["']([^"']+)["']`)
)

func (a *Animego) ParseEpisode(episode *extractors.Episode, player string, voicecover string) error {
	op := "animego.ParseEpisode"

	html, err := a.getEpisodeHTML(episode.ID)
	if err != nil {
		return err
	}

	links, err := a.parseLinks(html)
	if err != nil {
		log.Printf("error in pars links (%s): %v\n", op, err)
		return err
	}
	link := links[voicecover][player]
	if link == "" {
		log.Printf("cant find selected url for %s+%s\n", player, voicecover)
		return fmt.Errorf("cant find url for %s+%s", player, voicecover)
	}

	episode.PlayerURL = "https:" + link
	log.Printf("player url: %s", episode.PlayerURL)
	return nil
}

func (a *Animego) getEpisodeHTML(id string) ([]byte, error) {
	op := "animego.ParseEpisode.getEpisodeHTML"

	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("https://animego.me/player/videos/%s", id), nil)
	if err != nil {
		log.Printf("error in creating request (%s): %v\n", op, err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:148.0) Gecko/20100101 Firefox/148.0")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		log.Printf("error in get request (%s): %v\n", op, err)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("not ok status (%s): %d\n", op, resp.StatusCode)
		return nil, fmt.Errorf("not ok status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error in parsing response (%s), %v\n", op, err)
		return nil, err
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("error in unmarshal json (%s): %v\n", op, err)
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return []byte(result.Data.Content), nil
}

// ParseEpisodePlayers принимает HTML-фрагмент (или полную страницу) в виде []byte
// и возвращает PlayerLinks — map[озвучка]map[плеер]ссылка.
func (a *Animego) parseLinks(data []byte) (PlayerLinks, error) {
	result := make(PlayerLinks)

	for _, match := range reButton.FindAllSubmatch(data, -1) {
		attrs := match[1]

		playerURL := firstGroup(rePlayer, attrs)
		provider := firstGroup(reProvider, attrs)
		translation := firstGroup(reTranslation, attrs)

		if playerURL == "" || provider == "" || translation == "" {
			continue
		}

		playerURL = html.UnescapeString(playerURL)
		translation = html.UnescapeString(translation)
		provider = html.UnescapeString(provider)

		if result[translation] == nil {
			result[translation] = make(map[string]string)
		}
		if _, exists := result[translation][provider]; !exists {
			result[translation][provider] = playerURL
		}
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no player buttons found")
	}
	return result, nil
}

func firstGroup(re *regexp.Regexp, src []byte) string {
	m := re.FindSubmatch(src)
	if m == nil {
		return ""
	}
	return string(m[1])
}
