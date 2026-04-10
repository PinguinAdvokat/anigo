package yummyanime

import (
	"anigo/internal/extractors"
	"encoding/json"
	"log"
	"net/http"
	"slices"
	"sort"
	"strconv"
	"strings"
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

	temp := make(map[int]map[string]map[string]string)
	uniquePlayers := []string{}
	uniqueVoicecovers := []string{}
	for _, apiEpisode := range apiResponse.Episodes {
		num, _ := strconv.Atoi(apiEpisode.Number)
		player := strings.Replace(apiEpisode.Data.Player, "Плеер ", "", 1)
		voice := strings.Replace(apiEpisode.Data.Voicecover, "Озвучка ", "", 1)
		url := apiEpisode.IframeURL

		if !slices.Contains(uniquePlayers, player) && slices.Contains(extractors.AvailablePlayers, player) {
			uniquePlayers = append(uniquePlayers, player)
		}
		if !slices.Contains(uniqueVoicecovers, voice) {
			uniqueVoicecovers = append(uniqueVoicecovers, voice)
		}

		if !strings.HasPrefix(url, "https:") {
			url = "https:" + url
		}

		if temp[num] == nil {
			temp[num] = make(map[string]map[string]string)
		}
		if temp[num][player] == nil {
			temp[num][player] = make(map[string]string)
		}
		temp[num][player][voice] = url
	}

	// Шаг 2: Создание списка Episode
	var episodes []extractors.Episode
	episodeNums := make([]int, 0, len(temp))
	for num := range temp {
		episodeNums = append(episodeNums, num)
	}
	sort.Ints(episodeNums) // Сортировка по номеру серии

	for _, num := range episodeNums {
		episodes = append(episodes, extractors.Episode{
			AllVideos: temp[num],
		})
	}

	anime.Episodes = episodes
	anime.AvailablePlayers = uniquePlayers
	anime.AvailableVoiceover = uniqueVoicecovers

	return nil
}
