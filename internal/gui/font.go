package gui

import (
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
)

// FontConfig holds font configuration
type FontConfig struct {
	CustomFontPath string
	UsesSystemFont bool
}

// GetDefaultFontPaths returns a list of common font paths to try
func GetDefaultFontPaths() []string {
	return []string{
		// macOS system fonts (Chinese support)
		"/System/Library/Fonts/PingFang.ttc",
		"/System/Library/Fonts/STHeiti Light.ttc",
		"/System/Library/Fonts/Hiragino Sans GB.ttc",

		// User fonts on macOS
		filepath.Join(os.Getenv("HOME"), "Library/Fonts/PingFangSC-Regular.ttf"),

		// Common open-source fonts
		"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
		"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc",

		// Windows fonts (if running on Windows via Wine or dual boot)
		"C:\\Windows\\Fonts\\msyh.ttc",  // Microsoft YaHei
		"C:\\Windows\\Fonts\\simsun.ttc", // SimSun
	}
}

// FindAvailableFont searches for an available CJK font
func FindAvailableFont() (string, error) {
	paths := GetDefaultFontPaths()

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("no CJK font found in default paths")
}

// LoadCustomFont attempts to load a custom font file
func LoadCustomFont(fontPath string) error {
	if fontPath == "" {
		return fmt.Errorf("font path is empty")
	}

	// Check if file exists
	if _, err := os.Stat(fontPath); err != nil {
		return fmt.Errorf("font file not found: %s", fontPath)
	}

	// Set the FYNE_FONT environment variable
	os.Setenv("FYNE_FONT", fontPath)

	return nil
}

// SetupFont configures the font for the application
func SetupFont(customPath string) error {
	// If custom path is provided, use it
	if customPath != "" {
		return LoadCustomFont(customPath)
	}

	// Check if FYNE_FONT is already set
	if os.Getenv("FYNE_FONT") != "" {
		return nil // Already configured
	}

	// Try to find an available font
	fontPath, err := FindAvailableFont()
	if err != nil {
		// No font found, will use Fyne's default
		return fmt.Errorf("warning: %v, using default font", err)
	}

	return LoadCustomFont(fontPath)
}

// GetFontResource creates a font resource from a file path
func GetFontResource(fontPath string) (fyne.Resource, error) {
	if fontPath == "" {
		return nil, fmt.Errorf("font path is empty")
	}

	// Read font file
	data, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file: %w", err)
	}

	// Create resource
	fontName := filepath.Base(fontPath)
	resource := fyne.NewStaticResource(fontName, data)

	return resource, nil
}
