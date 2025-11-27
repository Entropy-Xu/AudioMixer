package main

import (
	"fmt"
	"os"

	"github.com/entropy/audio-mixer/internal/gui"
)

func main() {
	// Set Fyne to use system font for better CJK (Chinese/Japanese/Korean) support
	// This environment variable tells Fyne to use the system's font rendering
	if os.Getenv("FYNE_FONT") == "" {
		// On macOS, use PingFang SC for Chinese support
		if _, err := os.Stat("/System/Library/Fonts/PingFang.ttc"); err == nil {
			os.Setenv("FYNE_FONT", "/System/Library/Fonts/PingFang.ttc")
		}
	}

	app, err := gui.NewApp()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing application: %v\n", err)
		os.Exit(1)
	}

	app.Run()
}
