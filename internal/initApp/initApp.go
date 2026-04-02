package initApp

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func dataDir(appName string) string {
	if d := os.Getenv("XDG_DATA_HOME"); d != "" {
		return filepath.Join(d, appName)
	}
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Print("Failed get AppData directory\n")
		os.Exit(1)
	}
	return filepath.Join(home, ".local", "share", appName) // стандарт на Linux
}

func Init() string {
	appDir := dataDir("anigo")
	err := os.MkdirAll(appDir, 0755)
	if err != nil {
		log.Fatalf("Failed create directory: %v\n", err)
	}
	return appDir
}

func CreateLogFile(path string) *os.File {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("cant open log file", err)
	}
	return file
}
