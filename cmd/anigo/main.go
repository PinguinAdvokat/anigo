package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"
	"path/filepath"

	"anigo/internal/cache"
	"anigo/internal/initApp"
	"anigo/internal/parsers/kodik"
)

func main() {
	verbose := flag.Bool("v", false, "включить отладочные логи")
	flag.Parse()

	if !*verbose {
		// все log.Print, log.Printf и т.п. будут "тихими"
		log.SetOutput(io.Discard)
	}

	appDir := initApp.Init()
	cache := cache.New(filepath.Join(appDir, "cache.json"))

	// https://animego.me/anime/naruto-uragannye-hroniki-103

	client := &http.Client{}
	ctx := context.Background()

	parser := kodik.New(client, cache)
	parser.Parse(ctx, "https://kodik.info/seria/1567311/9ba0cd30f0ffa175bf8fc35b227f7ddc/720p?translations=false")
	parser.Parse(ctx, "https://kodik.info/seria/1570390/60887a958fb75529b9263abbe12889e9/720p?translations=false")
	parser.Parse(ctx, "https://kodik.info/seria/1573060/52341eea87588f673f9e8a6dbeadc2c2/720p?translations=false")
	parser.Parse(ctx, "https://kodik.info/seria/1567311/9ba0cd30f0ffa175bf8fc35b227f7ddc/720p?translations=false")
	cache.Save()
}
