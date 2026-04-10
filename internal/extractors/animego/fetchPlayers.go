package animego

import (
	"anigo/internal/extractors"
	"regexp"
	"slices"
)

type Response struct {
	Status string `json:"status"`
	Data   struct {
		Content string `json:"content"`
	} `json:"data"`
}

// FetchPlayers делает GET-запрос по url, парсит JSON-ответ
// и возвращает список уникальных плееров (например ["AniBoom", "Kodik", "Sibnet"]).
func (a *Animego) FetchPlayers(content string) ([]string, error) {
	// Ищем все вхождения data-provider-title="..." в HTML
	re := regexp.MustCompile(`data-provider-title="([^"]+)"`)
	matches := re.FindAllStringSubmatch(content, -1)

	// Дедуплицируем с сохранением порядка
	seen := make(map[string]bool)
	var players []string
	for _, m := range matches {
		name := m[1]
		if !seen[name] {
			seen[name] = true
			if slices.Contains(extractors.AvailablePlayers, name) {
				players = append(players, name)
			}
		}
	}

	return players, nil
}
