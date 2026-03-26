package animego

import (
	"anigo/internal/extractors"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// ldSchema is used to decode the JSON-LD block embedded in the page <head>.
type ldSchema struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	NumberEpisodes int    `json:"numberOfEpisodes"`
}

// ── Compiled regular expressions ─────────────────────────────────────────────

var (
	// JSON-LD structured data block.
	reLD = regexp.MustCompile(`<script type="application/ld\+json">([\s\S]*?)</script>`)

	// Dubbing link: href="/anime/dubbing/xxx" title="Name"
	reDubbing = regexp.MustCompile(`href="/anime/dubbing/[^"]*"\s+title="([^"]+)"`)

	// Player AJAX endpoint embedded in the page.
	rePlayerURL = regexp.MustCompile(`data-ajax-url="(/player/\d+)"`)

	// data-player attribute in the player response.
	rePlayerName = regexp.MustCompile(`data-player="([^"]+)"`)

	// Episode number label: data-label="25."
	reEpLabel = regexp.MustCompile(`data-label="(\d+)\."`)

	// Episode count from the info table (fallback).
	reEpCountGrid = regexp.MustCompile(`>Эпизоды</div>\s*<div [^>]*>(\d+)</div>`)

	// Strip HTML tags.
	reStripTags = regexp.MustCompile(`<[^>]+>`)
)

// ── Public API ───────────────────────────────────────────────────────────────

// ParseAnime fetches the anime page at pageURL, parses it, and returns AnimeInfo.
// It also fetches the separate /player/{id} endpoint to discover available players.
func (a *Animego) ParseAnime(anime extractors.Anime) (extractors.Anime, error) {
	pageURL := anime.URL
	body, err := a.fetchURL(pageURL)
	if err != nil {
		return extractors.Anime{}, fmt.Errorf("fetch page: %w", err)
	}

	anime = parseBody(anime, body)

	// fetch players
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(pageURL)

	if match == "" {
		err := fmt.Errorf("cant find anime id in %s", pageURL)
		log.Printf("error in getting players: %v\n", err)
		return extractors.Anime{}, err
	}
	num, _ := strconv.Atoi(match)
	playerContent, err := a.getPlayerContent(fmt.Sprintf("https://animego.me/player/%d", num))
	if err != nil {
		log.Printf("error in getting playerContent %v\n", err)
		return extractors.Anime{}, err
	}

	players, err := a.FetchPlayers(playerContent)
	if err != nil {
		log.Printf("error in geting players: %v\n", err)
		return extractors.Anime{}, err
	}
	anime.AvailablePlayers = players
	return anime, nil
}

// fetchURL performs an HTTP GET with browser-like headers.
func (a *Animego) fetchURL(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent",
		"Mozilla/5.0 (X11; Linux x86_64; rv:122.0) Gecko/20100101 Firefox/122.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,*/*;q=0.9")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en;q=0.8")
	req.Header.Set("Referer", "https://animego.me/")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, url)
	}

	return io.ReadAll(resp.Body)
}

// ── Internal helpers ─────────────────────────────────────────────────────────

// parseBody extracts all fields that are present in the static HTML.
func parseBody(anime extractors.Anime, body []byte) extractors.Anime {
	// ── JSON-LD (rich structured data) ────────────────────────────────────
	if m := reLD.FindSubmatch(body); len(m) == 2 {
		var ld ldSchema
		if err := json.Unmarshal(m[1], &ld); err == nil {
			anime.Title = ld.Name
			anime.Description = ld.Description
			anime.EpisodesCount = ld.NumberEpisodes
		}
	}

	// ── Episode count fallback ─────────────────────────────────────────────
	if anime.EpisodesCount == 0 {
		if m := reEpCountGrid.FindSubmatch(body); len(m) == 2 {
			anime.EpisodesCount, _ = strconv.Atoi(string(m[1]))
		}
	}

	// ── Voiceovers / Озвучки ──────────────────────────────────────────────
	seen := map[string]bool{}
	for _, m := range reDubbing.FindAllSubmatch(body, -1) {
		name := string(m[1])
		if !seen[name] {
			seen[name] = true
			anime.AvailableVoiceover = append(anime.AvailableVoiceover, name)
		}
	}
	return anime
}

// StripHTML removes HTML tags and collapses whitespace.
func StripHTML(s string) string {
	s = reStripTags.ReplaceAllString(s, " ")
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

// fetching https://animego.me/player/2422 and returning content string
func (a *Animego) getPlayerContent(url string) (*string, error) {
	req, err := http.NewRequest("GET", url, nil)
	log.Printf("url of request: %v\n", url)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Указываем заголовки — сервер может возвращать HTML без них
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("not ok status in fetchinf players: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body: %w", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("resp: %v\n", string(body))
		return nil, fmt.Errorf("unmarshal json: %w", err)
	}

	return &result.Data.Content, nil
}
