package kodik

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func (k *Kodik) getHTML(kodikURL string) (string, error) {
	op := "getPayload/getHTML"
	log.Printf("KodikUrl: %s", kodikURL)
	req, err := http.NewRequest(http.MethodGet, kodikURL, nil)
	if err != nil {
		log.Printf("error creating request in %v: %v\n", op, err)
		return "", err
	}

	resp, err := k.HttpClient.Do(req)
	if err != nil {
		log.Printf("error doing request in %v: %v\n", op, err)
		return "", err
	}

	defer resp.Body.Close()

	html, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error read resp in %v: %v\n", op, err)
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("get notOk status in %v: %v\n", op, err)
		return "", err
	}

	return string(html), nil
}

func (k *Kodik) getPayload(kodikURL string) (url.Values, error) {
	html, err := k.getHTML(kodikURL)
	if err != nil {
		return nil, err
	}
	payload := url.Values{}

	// Парсим type из vInfo.type = '...';
	typeRegex := regexp.MustCompile(`vInfo\.type\s*=\s*'([^']+)'`)
	if matches := typeRegex.FindStringSubmatch(html); len(matches) > 1 {
		payload.Set("type", matches[1])
	}

	// Парсим hash из vInfo.hash = '...';
	hashRegex := regexp.MustCompile(`vInfo\.hash\s*=\s*'([^']+)'`)
	if matches := hashRegex.FindStringSubmatch(html); len(matches) > 1 {
		payload.Set("hash", matches[1])
	}

	// Парсим id из vInfo.id = '...';
	idRegex := regexp.MustCompile(`vInfo\.id\s*=\s*'([^']+)'`)
	if matches := idRegex.FindStringSubmatch(html); len(matches) > 1 {
		payload.Set("id", matches[1])
	}

	return payload, nil
}
