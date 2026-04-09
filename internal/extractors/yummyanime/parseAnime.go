package yummyanime

import (
	"anigo/internal/extractors"
	"encoding/json"
	"log"
	"net/http"
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
	req, err := http.NewRequest("GET", anime.URL+"/videos", nil)
	if err != nil {
		log.Printf("YummyAnime ParseAnime request creation error: %v", err)
		return err
	}
	req.Header.Set("Host", "api.yani.tv")
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

	episodes := make(map[string]map[string][]extractors.Episode)

	maxIndex := 0
	for _, apiEpisode := range apiResponse.Episodes {
		player := apiEpisode.Data.Player
		voice := apiEpisode.Data.Voicecover

		// Инициализируем вложенную структуру если нужно
		if episodes[player] == nil {
			episodes[player] = make(map[string][]extractors.Episode)
		}
		if episodes[player][voice] == nil {
			episodes[player][voice] = []extractors.Episode{}
		}

		episodes[player][voice] = append(episodes[player][voice], extractors.Episode{
			PlayerURL: apiEpisode.IframeURL,
		})

		if len(episodes[player][voice]) > maxIndex {
			maxIndex = len(episodes[player][voice])
		}
	}
	anime.YummEpisodesRaw = episodes

	for _ = range maxIndex {
		anime.Episodes = append(anime.Episodes, extractors.Episode{})
	}

	uniquePlayers, uniqueVoicecovers := getUniqueLists(anime.YummEpisodesRaw)
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
