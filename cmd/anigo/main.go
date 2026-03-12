package main

import (
	"anigo/internal/app"
	"anigo/internal/cache"
	"anigo/internal/extractors/animego"
	"anigo/internal/initApp"
	"anigo/internal/manager"
	"anigo/internal/parsers/kodik"
	"flag"
	"io"
	"log"
	"net/http"
	"time"
)

func main() {
	verbose := flag.Bool("v", false, "включить отладочные логи")
	flag.Parse()

	if !*verbose {
		log.SetOutput(io.Discard)
	}

	appDir := initApp.Init()
	cache := cache.New(appDir)
	httpClient := &http.Client{Timeout: 5 * time.Second}
	kodikParser := kodik.New(httpClient, cache)
	animego := animego.New(kodikParser, httpClient)
	manager := manager.New(animego)

	app := app.New(manager)
	if err := app.EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
