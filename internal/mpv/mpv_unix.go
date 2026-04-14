//go:build !windows

package mpv

import (
	"fmt"
	"log"
	"net"
)

func New() *Mpv {
	return &Mpv{addr: "/tmp/mpv.sock"}
}

func (m *Mpv) Add(url string) error {
	if !m.isPlaying {
		err := fmt.Errorf("mpv not playing anything")
		log.Print(err)
		return err
	}
	conn, err := net.Dial("unix", "/tmp/mpv.sock")
	if err != nil {
		log.Printf("error in connection to socket: %v", err)
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(`{"command":["loadfile",` + `"` + url + `"` + `,"append"]}` + "\n"))
	if err != nil {
		log.Printf("error in writing to socket: %v", err)
		return err
	}
	return nil
}
