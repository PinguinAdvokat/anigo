package main

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"regexp"
)

func sanitizeFilename(name string) string {
	// Заменяем всё, что не буква/цифра/._- на "_"
	re := regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	name = re.ReplaceAllString(name, "_")
	if name == "" {
		name = "index"
	}
	return name + ".html"
}

func filenameFromURL(raw string) (string, error) {
	u, err := url.Parse(raw)
	if err != nil {
		return "", err
	}

	// Берём последний сегмент пути, если он пустой — "index"
	base := path.Base(u.Path)
	if base == "." || base == "/" || base == "" {
		base = "index"
	}

	// Добавим хост, чтобы имена были уникальнее
	name := u.Host + "_" + base
	return sanitizeFilename(name), nil
}

func downloadHTML(urlStr, filename string) error {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("error in create request: %w", err)
	}
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("content-encoding", "gzip")
	req.Header.Set("Referer", "https://animego.me/anime/van-pis-s1-65")
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("file create error: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, resp.Body); err != nil {
		return fmt.Errorf("copy error: %w", err)
	}

	return nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Использование: downloader <url>")
		return
	}

	rawURL := os.Args[1]

	filename, err := filenameFromURL(rawURL)
	if err != nil {
		fmt.Println("Ошибка парсинга URL:", err)
		return
	}

	if err := downloadHTML(rawURL, filename); err != nil {
		fmt.Println("Ошибка:", err)
		return
	}

	fmt.Printf("Страница %s сохранена в файл %s\n", rawURL, filename)
}
