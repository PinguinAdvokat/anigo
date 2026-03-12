package animego

import (
	"anigo/internal/extractors"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
)

type Response struct {
	Status string `json:"status"`
	Data   struct {
		Content string `json:"content"`
	} `json:"data"`
}

// FetchPlayers делает GET-запрос по url, парсит JSON-ответ
// и возвращает список уникальных плееров (например ["AniBoom", "Kodik", "Sibnet"]).
func (a *Animego) FetchPlayers(url string) ([]string, error) {
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

	// Ищем все вхождения data-provider-title="..." в HTML
	re := regexp.MustCompile(`data-provider-title="([^"]+)"`)
	matches := re.FindAllStringSubmatch(result.Data.Content, -1)

	// Дедуплицируем с сохранением порядка
	seen := make(map[string]bool)
	var players []string
	for _, m := range matches {
		name := m[1]
		if !seen[name] {
			seen[name] = true
			if contains(extractors.AvailablePlayers, name) {
				players = append(players, name)
			}
		}
	}

	return players, nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
