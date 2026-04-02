package kodik

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Response struct {
	Links LinkMap `json:"links"`
}

type LinkMap map[string][]LinkItem

type LinkItem struct {
	Src string `json:"src"`
}

// DecodeROT18AndBase64 декодирует строку применяя ROT18 (без поворота чисел), затем base64
func decodeURL(s string) (string, error) {
	// Сначала применяем ROT18
	decoded := rot18(s)

	// Затем декодируем из base64
	result, err := base64.RawURLEncoding.DecodeString(decoded)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// DecodeRot18 применяет к строке сдвиг ROT18 (каждый символ a-z/A-Z сдвигается на 18 позиций)
// Результат можно использовать как для кодирования, так и для декодирования (в отличие от ROT13,
// для восстановления исходной строки потребуется сдвиг на 8, но по условию требуется именно ROT18).
func rot18(s string) string {
	runes := []rune(s)
	for i, r := range runes {
		switch {
		case r >= 'a' && r <= 'z':
			runes[i] = 'a' + (r-'a'+18)%26
		case r >= 'A' && r <= 'Z':
			runes[i] = 'A' + (r-'A'+18)%26
		}
	}
	return string(runes)
}

func (k *Kodik) getSources(payload url.Values) (map[string]string, error) {
	op := "kodik/getSources"

	req, err := http.NewRequest(http.MethodPost, "https://kodikplayer.com/ftor", strings.NewReader(payload.Encode()))
	if err != nil {
		log.Printf("error creating request in %v: %v\n", op, err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(payload.Encode())))

	resp, err := k.HttpClient.Do(req)
	if err != nil {
		log.Printf("error doing request in %v: %v\n", op, err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("get notOk status in %v: %v\n", op, resp.Status)
		return nil, err
	}

	var parsed Response
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		log.Printf("error reading json response in %v: %v\n", op, err)
		return nil, err
	}

	sources := make(map[string]string)
	for qual, linkItem := range parsed.Links {
		url, err := decodeURL(linkItem[0].Src)
		if err != nil {
			log.Printf("error decodeUrl in %s: %v", op, err)
		} else {
			if strings.HasPrefix(url, "https:") {
				sources[qual] = url
			} else {
				sources[qual] = "https:" + url
			}
		}
	}
	log.Printf("sources: %v", sources)
	return sources, nil
}
