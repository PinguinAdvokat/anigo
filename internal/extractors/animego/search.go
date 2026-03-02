package animego

import (
	"anigo/internal/extractors"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (a *Animego) Search(query string) ([]extractors.AnimeInfo, error) {
	url := a.BaseURL + "/search/anime?q=" + query

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept-Language", "ru-RU,ru;q=0.9")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return a.parseAnimeList(resp.Body)
}

func (a *Animego) parseAnimeList(r io.Reader) ([]extractors.AnimeInfo, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var results []extractors.AnimeInfo

	doc.Find("div.ani-grid__item").Each(func(_ int, s *goquery.Selection) {
		// Название и URL из <a> внутри ani-grid__item-title
		titleLink := s.Find("div.ani-grid__item-title a").First()
		title := strings.TrimSpace(titleLink.Text())
		href, _ := titleLink.Attr("href")
		if href != "" && !strings.HasPrefix(href, "http") {
			href = a.BaseURL + href
		}

		// Рейтинг из div.rating-badge
		rating := strings.TrimSpace(s.Find("div.rating-badge").Text())

		if title == "" {
			return
		}

		results = append(results, extractors.AnimeInfo{
			Title:  title,
			Rating: rating,
			URL:    href,
		})
	})

	return results, nil
}
