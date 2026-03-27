package kodik

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type response struct {
	Links struct {
		R360 []struct {
			Src string `json:"Src"`
		} `json:"360"`
		R480 []struct {
			Src string `json:"Src"`
		} `json:"480"`
		R720 []struct {
			Src string `json:"Src"`
		} `json:"720"`
		R1080 []struct {
			Src string `json:"Src"`
		} `json:"1080"`
	} `json:"Links"`
}

// DecodeROT18AndBase64 декодирует строку применяя ROT18 (без поворота чисел), затем base64
func decodeURL(s string) (string, error) {
	// Сначала применяем ROT18
	decoded := rot18(s)

	// Затем декодируем из base64
	result, err := base64.StdEncoding.DecodeString(decoded)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

// rot18 применяет ROT18 к строке (только буквы, числа не меняются)
func rot18(s string) string {
	result := make([]rune, 0, len(s))

	for _, r := range s {
		if r >= 'a' && r <= 'z' {
			// Поворот на 18 позиций в нижнем регистре
			result = append(result, ((r-'a'+18)%26)+'a')
		} else if r >= 'A' && r <= 'Z' {
			// Поворот на 18 позиций в верхнем регистре
			result = append(result, ((r-'A'+18)%26)+'A')
		} else {
			// Числа и прочие символы не меняются
			result = append(result, r)
		}
	}

	return string(result)
}

func (k *Kodik) getSources(ctx context.Context, payload url.Values) (map[string]string, error) {
	op := "kodik/getSources"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://kodikplayer.com/ftor", strings.NewReader(payload.Encode()))
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

	var parsed response
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		log.Printf("error reading json response in %v: %v\n", op, err)
		return nil, err
	}

	sources := make(map[string]string)
	if parsed.Links.R360 != nil {
		url, err := decodeURL(parsed.Links.R360[0].Src)
		if err != nil {
			log.Printf("error decoding url in %v: %v\n", op, err)
			return nil, err
		}
		sources["360"] = "https:" + url
	}
	if parsed.Links.R480 != nil {
		url, err := decodeURL(parsed.Links.R480[0].Src)
		if err != nil {
			log.Printf("error decoding url in %v: %v\n", op, err)
			return nil, err
		}
		sources["480"] = "https:" + url
	}
	if parsed.Links.R720 != nil {
		url, err := decodeURL(parsed.Links.R720[0].Src)
		if err != nil {
			log.Printf("error decoding url in %v: %v\n", op, err)
			return nil, err
		}
		sources["720"] = "https:" + url
	}
	if parsed.Links.R1080 != nil {
		url, err := decodeURL(parsed.Links.R1080[0].Src)
		if err != nil {
			log.Printf("error decoding url in %v: %v\n", op, err)
			return nil, err
		}
		sources["1080"] = "https:" + url
	}
	return sources, nil
}
