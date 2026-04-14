package mpv

import (
	"fmt"
	"io"
	"log"
	"os/exec"
)

func (m *Mpv) Play(url string) error {
	if m.isPlaying {
		m.Add(url)
		return fmt.Errorf("mpv player already playing video")
	}
	cmd := exec.Command(
		"mpv",
		"--save-position-on-quit",
		fmt.Sprintf("--input-ipc-server=%s", m.addr),
		url,
	)

	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard

	if err := cmd.Start(); err != nil {
		return err
	}
	m.isPlaying = true
	go func() {
		err := cmd.Wait()
		m.isPlaying = false
		if err != nil {
			log.Printf("error in mpv: %v", err)
			return
		}
	}()
	return nil
}
