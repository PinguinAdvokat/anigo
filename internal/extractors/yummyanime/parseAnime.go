package yummyanime

import (
	"anigo/internal/extractors"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type ApiEpisode struct {
	Number    string `json:"number"`
	IframeURL string `json:"iframe_url"`
	VideoID   int    `json:"video_id"`
	Data      data   `json:"data"`
}

type data struct {
	Player     string `json:"player"`
	Voicecover string `json:"dubbing"`
	PlayerID   int    `json:"player_id"`
}

func (y *YummyAnime) ParseAnime(anime *extractors.Anime) error {
	req, err := http.NewRequest("GET", anime.URL, nil)
	if err != nil {
		log.Printf("YummyAnime ParseAnime request creation error: %v", err)
		return err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Set("Accept-Language", "ru")

	resp, err := y.httpClient.Do(req)
	if err != nil {
		log.Printf("YummyAnime ParseAnime request error: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("YummyAnime ParseAnime HTTP error: %d", resp.StatusCode)
		return err
	}

	var apiResponse struct {
		Episodes []ApiEpisode `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&apiResponse); err != nil {
		log.Printf("YummyAnime ParseAnime response decode error: %v", err)
		return err
	}

	for _, apiEpisode := range apiResponse.Episodes {
		number, err := strconv.Atoi(apiEpisode.Number)
		if err != nil {
			log.Printf("error in setting yummy episodes: %v", err)
			continue
		}
		anime.YummyApisodesRaw[apiEpisode.Data.Player][apiEpisode.Data.Voicecover][number-1] = extractors.Episode{
			PlayerURL: apiEpisode.IframeURL,
		}
	}

	uniquePlayers, uniqueVoicecovers := getUniqueLists(anime.YummyApisodesRaw)
	anime.AvailablePlayers = uniquePlayers
	anime.AvailableVoiceover = uniqueVoicecovers

	return nil
}

func getUniqueLists[K1 comparable, K2 comparable, V any](
	data map[K1]map[K2]V,
) ([]K1, []K2) {
	keysOuter := make([]K1, 0)
	keysInner := make(map[K2]struct{})

	// уникальные внешние ключи
	for kOuter := range data {
		keysOuter = append(keysOuter, kOuter)
	}

	// уникальные внутренние ключи
	for _, innerMap := range data {
		for kInner := range innerMap {
			keysInner[kInner] = struct{}{}
		}
	}

	// переводим внутренние ключи в срез
	keysInnerSlice := make([]K2, 0, len(keysInner))
	for k := range keysInner {
		keysInnerSlice = append(keysInnerSlice, k)
	}

	return keysOuter, keysInnerSlice
}
