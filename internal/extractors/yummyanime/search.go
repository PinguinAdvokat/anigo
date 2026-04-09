package yummyanime

import (
	"anigo/internal/extractors"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	req, err := http.NewRequest("GET", y.BaseURL+"/search?q="+name, nil)
	if err != nil {
		log.Printf("YummyAnime search request creation error: %v", err)
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept-Language", "ru")

	resp, err := y.httpClient.Do(req)
	if err != nil {
		log.Printf("YummyAnime search request error: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("YummyAnime search HTTP error: %d", resp.StatusCode)
		return nil, err
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
