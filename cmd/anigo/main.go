package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/PinguinAdvokat/anigo/internal/app"
	"github.com/PinguinAdvokat/anigo/internal/cache"
	"github.com/PinguinAdvokat/anigo/internal/initApp"
	"github.com/PinguinAdvokat/anigo/internal/manager"
	"github.com/PinguinAdvokat/anigo/internal/mpv"
	"github.com/PinguinAdvokat/anigo/internal/parsers/kodik"
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
	manager := manager.New("", kodikParser, httpClient)
	mpv := mpv.New()

	app := app.New(manager, mpv, httpClient)
	if err := app.EnableMouse(true).EnablePaste(true).Run(); err != nil {
		panic(err)
	}
}
