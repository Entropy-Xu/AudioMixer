//go:build !cgo
// +build !cgo

package audio

import (
	"fmt"
)

// DeviceInfo holds information about an audio device
type DeviceInfo struct {
	Index             int
	Name              string
	MaxInputChannels  int
	MaxOutputChannels int
	DefaultSampleRate float64
	IsDefaultInput    bool
	IsDefaultOutput   bool
	HostAPI           string
}

// DeviceManager handles audio device enumeration and management
type DeviceManager struct {
	initialized bool
}

// NewDeviceManager creates a new device manager
func NewDeviceManager() *DeviceManager {
	return &DeviceManager{}
}

// Initialize initializes PortAudio (stub - not available without CGO)
func (dm *DeviceManager) Initialize() error {
	dm.initialized = true
	return nil
}

// Terminate cleans up PortAudio (stub)
func (dm *DeviceManager) Terminate() error {
	dm.initialized = false
	return nil
}

// ListDevices returns all available audio devices (stub)
func (dm *DeviceManager) ListDevices() ([]*DeviceInfo, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	// Return empty list when CGO is disabled
	// On Windows, WASAPI will be used instead
	return []*DeviceInfo{}, nil
}

// GetInputDevices returns only input-capable devices (stub)
func (dm *DeviceManager) GetInputDevices() ([]*DeviceInfo, error) {
	return dm.ListDevices()
}

// GetOutputDevices returns only output-capable devices (stub)
func (dm *DeviceManager) GetOutputDevices() ([]*DeviceInfo, error) {
	return dm.ListDevices()
}
