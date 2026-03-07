package main

import (
	"anigo/internal/app"
	"flag"
	"io"
	"log"
)

func main() {
	verbose := flag.Bool("v", false, "включить отладочные логи")
	flag.Parse()

	if !*verbose {
		// все log.Print, log.Printf и т.п. будут "тихими"
		log.SetOutput(io.Discard)
	}

	app := app.New()
	if err := app.EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
