package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/entropy/audio-mixer/internal/gui"
)

func main() {
	// Command line flags
	fontPath := flag.String("font", "", "Path to custom font file (TTF/TTC) for better CJK support")
	flag.Parse()

	// Setup font for CJK (Chinese/Japanese/Korean) support
	if err := gui.SetupFont(*fontPath); err != nil {
		// Non-fatal error, just print warning
		fmt.Fprintf(os.Stderr, "Font setup: %v\n", err)
	}

	app, err := gui.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	app.Run()
}
