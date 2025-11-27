package main

import (
	"fmt"
	"os"

	"github.com/entropy/audio-mixer/internal/gui"
)

func main() {
	app, err := gui.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	app.Run()
}
