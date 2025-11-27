package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the application configuration
type Config struct {
	// Audio settings
	SampleRate   float64 `json:"sample_rate"`
	BufferSize   int     `json:"buffer_size"`
	Channels     int     `json:"channels"`

	// Device indices
	Input1DeviceIndex int `json:"input1_device_index"` // Microphone/Line input
	Input2DeviceIndex int `json:"input2_device_index"` // System audio (loopback device)
	OutputDeviceIndex int `json:"output_device_index"` // Virtual output (BlackHole, etc.)

	// Virtual device settings
	UseVirtualOutput bool   `json:"use_virtual_output"` // Use virtual device for output
	LoopbackDeviceName string `json:"loopback_device_name"` // Name of loopback device

	// Volume settings (0.0 to 2.0)
	Input1Gain float32 `json:"input1_gain"`
	Input2Gain float32 `json:"input2_gain"`
	MasterGain float32 `json:"master_gain"`

	// UI preferences
	WindowWidth  int  `json:"window_width"`
	WindowHeight int  `json:"window_height"`
	StartMinimized bool `json:"start_minimized"`
}

// DefaultConfig returns a configuration with default values
func DefaultConfig() *Config {
	return &Config{
		SampleRate:         48000,
		BufferSize:         512,
		Channels:           2,
		Input1DeviceIndex:  -1, // -1 means use default device
		Input2DeviceIndex:  -1, // Will auto-detect loopback device
		OutputDeviceIndex:  -1, // Will auto-detect virtual output
		UseVirtualOutput:   true, // Use virtual device by default
		LoopbackDeviceName: "BlackHole", // Default to BlackHole on macOS
		Input1Gain:         1.0,
		Input2Gain:         1.0,
		MasterGain:         1.0,
		WindowWidth:        800,
		WindowHeight:       600,
		StartMinimized:     false,
	}
}

// ConfigManager handles loading and saving configuration
type ConfigManager struct {
	configPath string
}

// NewConfigManager creates a new configuration manager
func NewConfigManager() (*ConfigManager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".audio-mixer")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "config.json")

	return &ConfigManager{
		configPath: configPath,
	}, nil
}

// Load loads configuration from disk, returns default config if file doesn't exist
func (cm *ConfigManager) Load() (*Config, error) {
	// Check if config file exists
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Return default config if file doesn't exist
		return DefaultConfig(), nil
	}

	// Read config file
	data, err := os.ReadFile(cm.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Parse JSON
	config := DefaultConfig()
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return config, nil
}

// Save saves configuration to disk
func (cm *ConfigManager) Save(config *Config) error {
	// Validate configuration
	if err := cm.validateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration: %w", err)
	}

	// Marshal to JSON with indentation
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(cm.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// validateConfig validates configuration values
func (cm *ConfigManager) validateConfig(config *Config) error {
	if config.SampleRate <= 0 {
		return fmt.Errorf("sample rate must be positive")
	}

	if config.BufferSize <= 0 {
		return fmt.Errorf("buffer size must be positive")
	}

	if config.Channels <= 0 || config.Channels > 2 {
		return fmt.Errorf("channels must be 1 or 2")
	}

	if config.Input1Gain < 0 || config.Input1Gain > 2.0 {
		return fmt.Errorf("input1 gain must be between 0.0 and 2.0")
	}

	if config.Input2Gain < 0 || config.Input2Gain > 2.0 {
		return fmt.Errorf("input2 gain must be between 0.0 and 2.0")
	}

	if config.MasterGain < 0 || config.MasterGain > 2.0 {
		return fmt.Errorf("master gain must be between 0.0 and 2.0")
	}

	return nil
}

// GetConfigPath returns the path to the configuration file
func (cm *ConfigManager) GetConfigPath() string {
	return cm.configPath
}
