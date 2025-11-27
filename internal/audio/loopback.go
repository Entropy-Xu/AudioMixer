package audio

import (
	"fmt"
	"runtime"

	"github.com/gordonklaus/portaudio"
)

// LoopbackDevice represents a virtual loopback audio device
type LoopbackDevice struct {
	Name   string
	Device *portaudio.DeviceInfo
}

// FindLoopbackDevice finds a suitable loopback/virtual device for output
// On macOS, this looks for BlackHole or similar virtual audio drivers
// On Windows, this looks for VB-Cable or similar
// On Linux, this looks for PulseAudio null sink or ALSA loopback
func FindLoopbackDevice(dm *DeviceManager) (*LoopbackDevice, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	// Platform-specific device names to search for
	var searchNames []string
	switch runtime.GOOS {
	case "darwin": // macOS
		searchNames = []string{
			"BlackHole",
			"Soundflower",
			"Loopback Audio",
		}
	case "windows":
		searchNames = []string{
			"CABLE Input",
			"VB-Audio",
			"Virtual Audio Cable",
		}
	case "linux":
		searchNames = []string{
			"pulse",
			"null",
			"Loopback",
		}
	default:
		return nil, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	// Search for virtual audio devices
	for _, dev := range devices {
		deviceName := dev.Name

		// Check if device supports both input and output (loopback characteristic)
		if dev.MaxInputChannels > 0 && dev.MaxOutputChannels > 0 {
			for _, searchName := range searchNames {
				if contains(deviceName, searchName) {
					return &LoopbackDevice{
						Name:   deviceName,
						Device: dev,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("no loopback device found. Please install: macOS(BlackHole), Windows(VB-Cable), Linux(PulseAudio loopback)")
}

// GetLoopbackDeviceByName finds a loopback device by exact or partial name match
func GetLoopbackDeviceByName(dm *DeviceManager, name string) (*LoopbackDevice, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	for _, dev := range devices {
		if contains(dev.Name, name) && dev.MaxInputChannels > 0 && dev.MaxOutputChannels > 0 {
			return &LoopbackDevice{
				Name:   dev.Name,
				Device: dev,
			}, nil
		}
	}

	return nil, fmt.Errorf("loopback device '%s' not found", name)
}

// ListLoopbackDevices returns all detected loopback/virtual devices
func ListLoopbackDevices(dm *DeviceManager) ([]*LoopbackDevice, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	var loopbackDevices []*LoopbackDevice

	// Platform-specific keywords
	var keywords []string
	switch runtime.GOOS {
	case "darwin":
		keywords = []string{"BlackHole", "Soundflower", "Loopback"}
	case "windows":
		keywords = []string{"CABLE", "VB-Audio", "Virtual"}
	case "linux":
		keywords = []string{"pulse", "null", "Loopback", "Monitor"}
	}

	for _, dev := range devices {
		// Virtual devices typically support both input and output
		if dev.MaxInputChannels > 0 && dev.MaxOutputChannels > 0 {
			for _, keyword := range keywords {
				if contains(dev.Name, keyword) {
					loopbackDevices = append(loopbackDevices, &LoopbackDevice{
						Name:   dev.Name,
						Device: dev,
					})
					break
				}
			}
		}
	}

	return loopbackDevices, nil
}

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	// Simple case-insensitive contains
	sLower := toLower(s)
	substrLower := toLower(substr)

	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

// toLower converts a string to lowercase
func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}
