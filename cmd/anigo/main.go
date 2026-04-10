package main

import (
	"anigo/internal/app"
	"anigo/internal/cache"
	"anigo/internal/extractors/yummyanime"
	"anigo/internal/initApp"
	"anigo/internal/manager"
	"anigo/internal/mpv"
	"anigo/internal/parsers/kodik"
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

func main() {
	logFileFlag := flag.String("logfile", "", "path to file for log writing (creating if not exist)")
	flag.Parse()

	appDir := initApp.Init()

	logFilePath := *logFileFlag
	if logFilePath == "" {
		logFilePath = filepath.Join(appDir, "anigo.log")
	}
	logFile := initApp.CreateLogFile(logFilePath)
	defer logFile.Close()
	log.SetOutput(logFile)

	cache := cache.New(appDir)
	httpClient := &http.Client{Timeout: 3 * time.Second}
	kodikParser := kodik.New(httpClient, cache)
	// animego := animego.New(httpClient)
	yummyanime := yummyanime.New(httpClient)
	mpv := mpv.New()
	manager := manager.New(yummyanime, kodikParser)

	app := app.New(manager, mpv, httpClient)
	if err := app.EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
