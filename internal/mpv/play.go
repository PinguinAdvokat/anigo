package mpv

import (
	"os"
	"os/exec"
)

func (m *Mpv) Play(watchLaterDir string, url string) error {
	// RAFACTOR THIS BEFORE USE!!!!!!!!!
	cmd := exec.Command(
		"mpv",
		"--save-position-on-quit",
		"--watch-later-directory="+watchLaterDir,
		url,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return err
	}
	return nil
}
