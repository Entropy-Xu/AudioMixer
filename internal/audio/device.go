package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

// DeviceInfo holds information about an audio device
type DeviceInfo struct {
	Index              int
	Name               string
	MaxInputChannels   int
	MaxOutputChannels  int
	DefaultSampleRate  float64
	IsDefaultInput     bool
	IsDefaultOutput    bool
	HostAPI            string
}

// DeviceManager handles audio device enumeration and management
type DeviceManager struct {
	initialized bool
}

// NewDeviceManager creates a new device manager
func NewDeviceManager() *DeviceManager {
	return &DeviceManager{}
}

// Initialize initializes PortAudio
func (dm *DeviceManager) Initialize() error {
	if dm.initialized {
		return nil
	}

	if err := portaudio.Initialize(); err != nil {
		return fmt.Errorf("failed to initialize PortAudio: %w", err)
	}

	dm.initialized = true
	return nil
}

// Terminate cleans up PortAudio
func (dm *DeviceManager) Terminate() error {
	if !dm.initialized {
		return nil
	}

	if err := portaudio.Terminate(); err != nil {
		return fmt.Errorf("failed to terminate PortAudio: %w", err)
	}

	dm.initialized = false
	return nil
}

// ListDevices returns all available audio devices
func (dm *DeviceManager) ListDevices() ([]*DeviceInfo, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	defaultInput, _ := portaudio.DefaultInputDevice()
	defaultOutput, _ := portaudio.DefaultOutputDevice()

	var deviceList []*DeviceInfo
	for i, dev := range devices {
		hostAPI, _ := portaudio.HostApi(dev.HostApi)
		hostAPIName := "Unknown"
		if hostAPI != nil {
			hostAPIName = hostAPI.Name
		}

		info := &DeviceInfo{
			Index:              i,
			Name:               dev.Name,
			MaxInputChannels:   dev.MaxInputChannels,
			MaxOutputChannels:  dev.MaxOutputChannels,
			DefaultSampleRate:  dev.DefaultSampleRate,
			IsDefaultInput:     defaultInput != nil && dev == defaultInput,
			IsDefaultOutput:    defaultOutput != nil && dev == defaultOutput,
			HostAPI:            hostAPIName,
		}
		deviceList = append(deviceList, info)
	}

	return deviceList, nil
}

// GetInputDevices returns only input-capable devices
func (dm *DeviceManager) GetInputDevices() ([]*DeviceInfo, error) {
	allDevices, err := dm.ListDevices()
	if err != nil {
		return nil, err
	}

	var inputDevices []*DeviceInfo
	for _, dev := range allDevices {
		if dev.MaxInputChannels > 0 {
			inputDevices = append(inputDevices, dev)
		}
	}

	return inputDevices, nil
}

// GetOutputDevices returns only output-capable devices
func (dm *DeviceManager) GetOutputDevices() ([]*DeviceInfo, error) {
	allDevices, err := dm.ListDevices()
	if err != nil {
		return nil, err
	}

	var outputDevices []*DeviceInfo
	for _, dev := range allDevices {
		if dev.MaxOutputChannels > 0 {
			outputDevices = append(outputDevices, dev)
		}
	}

	return outputDevices, nil
}

// GetDeviceByIndex returns a specific device by index
func (dm *DeviceManager) GetDeviceByIndex(index int) (*portaudio.DeviceInfo, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	devices, err := portaudio.Devices()
	if err != nil {
		return nil, fmt.Errorf("failed to enumerate devices: %w", err)
	}

	if index < 0 || index >= len(devices) {
		return nil, fmt.Errorf("device index %d out of range", index)
	}

	return devices[index], nil
}

// GetDefaultInputDevice returns the default input device
func (dm *DeviceManager) GetDefaultInputDevice() (*portaudio.DeviceInfo, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	dev, err := portaudio.DefaultInputDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get default input device: %w", err)
	}

	return dev, nil
}

// GetDefaultOutputDevice returns the default output device
func (dm *DeviceManager) GetDefaultOutputDevice() (*portaudio.DeviceInfo, error) {
	if !dm.initialized {
		return nil, fmt.Errorf("device manager not initialized")
	}

	dev, err := portaudio.DefaultOutputDevice()
	if err != nil {
		return nil, fmt.Errorf("failed to get default output device: %w", err)
	}

	return dev, nil
}
