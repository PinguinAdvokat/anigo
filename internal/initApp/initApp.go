package initApp

import (
	"fmt"
	"os"
	"os/exec"
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
	_, err := exec.LookPath("mpv")
	if err != nil {
		fmt.Print("Cant find mpv player\n")
		//os.Exit(0)
	}
	appDir := dataDir("anigo")
	err = os.MkdirAll(appDir, 0755)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed create directory: %v\n", err)
		os.Exit(1)
	}

	return appDir
}
