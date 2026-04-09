package yummyanime

import (
	"anigo/internal/extractors"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
)

type ApiAnime struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	ID          int    `json:"anime_id"`
	CoverURL    struct {
		Medium string `json:"medium"`
	} `json:"poster"`
	Rating struct {
		Shikimori float64 `json:"shikimori_rating"`
	} `json:"rating"`
}

func (y *YummyAnime) Search(name string) ([]extractors.Anime, error) {
	encoded := url.QueryEscape(name)
	url := fmt.Sprintf("%s/search?q=%s", y.BaseURL, encoded)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("YummyAnime search request creation error: %v", err)
		return nil, err
	}
	req.Host = "api.yani.tv"
	req.Header.Set("User-Agent", "PostmanRuntime/7.49.1")
	req.Header.Set("Acept-Language", "ru")

	resp, err := y.httpClient.Do(req)
	if err != nil {
		log.Printf("YummyAnime search request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("YummyAnime search HTTP error: %d", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		log.Printf("returned body: %v", string(body))
		if err != nil {
			log.Print(err)
		}
		return nil, fmt.Errorf("YummyAnime search HTTP error: %d", resp.StatusCode)
	}

	var apiResponse struct {
		Results []ApiAnime `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("YummyAnime search response decode error: %v", err)
		return nil, err
	}

	var results []extractors.Anime
	for _, apiAnime := range apiResponse.Results {
		results = append(results, extractors.Anime{
			URL:         fmt.Sprintf("%s/anime/%d", y.BaseURL, apiAnime.ID),
			Title:       apiAnime.Title,
			Description: apiAnime.Description,
			CoverURL:    "https:" + apiAnime.CoverURL.Medium,
			Rating:      fmt.Sprintf("%.1f", apiAnime.Rating.Shikimori),
		})
	}

	return results, nil
}
