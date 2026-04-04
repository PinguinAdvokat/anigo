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

	// Cover image URL from preload link: href="https://img.cdngos.com/anime/..."
	reCoverURL = regexp.MustCompile(`<link\s+rel="preload"\s+as="image"[^>]*href="([^"]+)"`)
)

// ── Public API ───────────────────────────────────────────────────────────────

// ParseAnime fetches the anime page at pageURL, parses it, and returns AnimeInfo.
// It also fetches the separate /player/{id} endpoint to discover available players.
func (a *Animego) ParseAnime(anime *extractors.Anime) error {
	body, err := a.fetchURL(anime.URL)
	if err != nil {
		return fmt.Errorf("fetch page: %w", err)
	}

	// base anime info (voicecovers, players, episodes count)
	parseBody(anime, body)
	log.Printf("Cover URL %s", anime.CoverURL)

	// get anime id
	re := regexp.MustCompile(`\d+$`)
	match := re.FindString(anime.URL)
	if match == "" {
		err := fmt.Errorf("cant find anime id in %s", anime.URL)
		log.Printf("error in getting players: %v\n", err)
		return err
	}
	num, _ := strconv.Atoi(match)

	// get player html from content in json
	playerContent, err := a.getPlayerContent(fmt.Sprintf("https://animego.me/player/%d", num))
	if err != nil {
		log.Printf("error in getting playerContent %v\n", err)
		return err
	}

	// parsing players from html
	players, err := a.FetchPlayers(playerContent)
	if err != nil {
		log.Printf("error in geting players: %v\n", err)
		return err
	}
	anime.AvailablePlayers = players

	// parsing episodes (name, id)
	anime.Episodes = a.parseEpisodes(playerContent)

	return nil
}

func (a *Animego) parseEpisodes(html string) []extractors.Episode {
	re := regexp.MustCompile(`data-episode-title="([^"]*)"[^>]*data-episode="(\d+)"`)
	matches := re.FindAllStringSubmatch(html, -1)

	episodes := make([]extractors.Episode, 0, len(matches))
	for _, m := range matches {
		id, err := strconv.Atoi(m[2])
		if err != nil {
			continue
		}
		episodes = append(episodes, extractors.Episode{
			Title: m[1],
			ID:    strconv.Itoa(id),
		})
	}
	return episodes
}

// getCoverURL extracts the anime cover image URL from the HTML page.
// It looks for a preload link with rel="preload" and as="image".
func getCoverURL(anime *extractors.Anime, body []byte) {
	if m := reCoverURL.FindSubmatch(body); len(m) == 2 {
		anime.CoverURL = string(m[1])
		return
	}

	// Fallback: try to extract from JSON-LD structured data
	if m := reLD.FindSubmatch(body); len(m) == 2 {
		var ld ldSchema
		if err := json.Unmarshal(m[1], &ld); err == nil && ld.Description != "" {
			// Extract image URL from JSON-LD if preload link not found
			type ldFull struct {
				Image string `json:"image"`
			}
			var ldFullData ldFull
			if err := json.Unmarshal(m[1], &ldFullData); err == nil {
				anime.CoverURL = ldFullData.Image
			}
		}
	}
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
func parseBody(anime *extractors.Anime, body []byte) {
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

	// ── Cover image URL ───────────────────────────────────────────────────
	getCoverURL(anime, body)
}

// StripHTML removes HTML tags and collapses whitespace.
func StripHTML(s string) string {
	s = reStripTags.ReplaceAllString(s, " ")
	return strings.TrimSpace(strings.Join(strings.Fields(s), " "))
}

// fetching https://animego.me/player/2422 and returning content string
func (a *Animego) getPlayerContent(url string) (string, error) {
	req, err := http.NewRequest("GET", url, nil)
	log.Printf("url of request: %v\n", url)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	// Указываем заголовки — сервер может возвращать HTML без них
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("http get: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("not ok status in fetchinfo players: %d\n", resp.StatusCode)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read body: %w", err)
	}

	var result Response
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("resp: %v\n", string(body))
		return "", fmt.Errorf("unmarshal json: %w", err)
	}

	return result.Data.Content, nil
}
