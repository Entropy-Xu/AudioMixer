package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/entropy/audio-mixer/internal/audio"
	"github.com/entropy/audio-mixer/internal/config"
)

func main() {
	fmt.Println("=== Audio Mixer ===")
	fmt.Println("Cross-platform audio mixing tool")
	fmt.Println()

	// Initialize configuration manager
	configManager, err := config.NewConfigManager()
	if err != nil {
		fmt.Printf("Error initializing config manager: %v\n", err)
		os.Exit(1)
	}

	// Load configuration
	cfg, err := configManager.Load()
	if err != nil {
		fmt.Printf("Error loading config: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Configuration loaded from: %s\n\n", configManager.GetConfigPath())

	// Initialize device manager
	deviceManager := audio.NewDeviceManager()
	if err := deviceManager.Initialize(); err != nil {
		fmt.Printf("Error initializing PortAudio: %v\n", err)
		os.Exit(1)
	}
	defer deviceManager.Terminate()

	// List available devices
	fmt.Println("Available Audio Devices:")
	fmt.Println("------------------------")

	inputDevices, err := deviceManager.GetInputDevices()
	if err != nil {
		fmt.Printf("Error listing input devices: %v\n", err)
		os.Exit(1)
	}

	outputDevices, err := deviceManager.GetOutputDevices()
	if err != nil {
		fmt.Printf("Error listing output devices: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\nInput Devices:")
	for _, dev := range inputDevices {
		defaultMarker := ""
		if dev.IsDefaultInput {
			defaultMarker = " [DEFAULT]"
		}
		fmt.Printf("  [%d] %s (Channels: %d, SR: %.0f Hz, API: %s)%s\n",
			dev.Index, dev.Name, dev.MaxInputChannels, dev.DefaultSampleRate, dev.HostAPI, defaultMarker)
	}

	fmt.Println("\nOutput Devices:")
	for _, dev := range outputDevices {
		defaultMarker := ""
		if dev.IsDefaultOutput {
			defaultMarker = " [DEFAULT]"
		}
		fmt.Printf("  [%d] %s (Channels: %d, SR: %.0f Hz, API: %s)%s\n",
			dev.Index, dev.Name, dev.MaxOutputChannels, dev.DefaultSampleRate, dev.HostAPI, defaultMarker)
	}

	// Interactive device selection
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\n=== Device Configuration ===")

	// Select Input 1 (Microphone)
	fmt.Printf("\nSelect Input 1 device (Microphone) [current: %d, -1 for default]: ", cfg.Input1DeviceIndex)
	input1Idx := readInt(reader, cfg.Input1DeviceIndex)
	cfg.Input1DeviceIndex = input1Idx

	// Select Input 2 (Application Audio)
	fmt.Printf("Select Input 2 device (Application Audio) [current: %d, -1 for default, -2 to skip]: ", cfg.Input2DeviceIndex)
	input2Idx := readInt(reader, cfg.Input2DeviceIndex)
	cfg.Input2DeviceIndex = input2Idx

	// Select Output
	fmt.Printf("Select Output device [current: %d, -1 for default]: ", cfg.OutputDeviceIndex)
	outputIdx := readInt(reader, cfg.OutputDeviceIndex)
	cfg.OutputDeviceIndex = outputIdx

	// Volume settings
	fmt.Println("\n=== Volume Configuration (0.0 - 2.0) ===")
	fmt.Printf("Input 1 Gain [current: %.2f]: ", cfg.Input1Gain)
	input1Gain := readFloat32(reader, cfg.Input1Gain)
	cfg.Input1Gain = input1Gain

	fmt.Printf("Input 2 Gain [current: %.2f]: ", cfg.Input2Gain)
	input2Gain := readFloat32(reader, cfg.Input2Gain)
	cfg.Input2Gain = input2Gain

	fmt.Printf("Master Gain [current: %.2f]: ", cfg.MasterGain)
	masterGain := readFloat32(reader, cfg.MasterGain)
	cfg.MasterGain = masterGain

	// Save configuration
	if err := configManager.Save(cfg); err != nil {
		fmt.Printf("Warning: Failed to save config: %v\n", err)
	}

	// Setup mixer configuration
	mixerConfig := audio.DefaultMixerConfig()
	mixerConfig.SampleRate = cfg.SampleRate
	mixerConfig.BufferSize = cfg.BufferSize
	mixerConfig.Channels = cfg.Channels
	mixerConfig.Input1Gain = cfg.Input1Gain
	mixerConfig.Input2Gain = cfg.Input2Gain
	mixerConfig.MasterGain = cfg.MasterGain

	// Get device info
	if cfg.Input1DeviceIndex >= 0 {
		dev, err := deviceManager.GetDeviceByIndex(cfg.Input1DeviceIndex)
		if err != nil {
			fmt.Printf("Error getting input1 device: %v\n", err)
			os.Exit(1)
		}
		mixerConfig.Input1Device = dev
	} else if cfg.Input1DeviceIndex == -1 {
		dev, err := deviceManager.GetDefaultInputDevice()
		if err != nil {
			fmt.Printf("Error getting default input device: %v\n", err)
			os.Exit(1)
		}
		mixerConfig.Input1Device = dev
	}

	if cfg.Input2DeviceIndex >= 0 {
		dev, err := deviceManager.GetDeviceByIndex(cfg.Input2DeviceIndex)
		if err != nil {
			fmt.Printf("Error getting input2 device: %v\n", err)
			os.Exit(1)
		}
		mixerConfig.Input2Device = dev
	} else if cfg.Input2DeviceIndex == -1 {
		dev, err := deviceManager.GetDefaultInputDevice()
		if err != nil {
			fmt.Printf("Warning: Could not get default input device for input2: %v\n", err)
		} else {
			mixerConfig.Input2Device = dev
		}
	}

	if cfg.OutputDeviceIndex >= 0 {
		dev, err := deviceManager.GetDeviceByIndex(cfg.OutputDeviceIndex)
		if err != nil {
			fmt.Printf("Error getting output device: %v\n", err)
			os.Exit(1)
		}
		mixerConfig.OutputDevice = dev
	} else if cfg.OutputDeviceIndex == -1 {
		dev, err := deviceManager.GetDefaultOutputDevice()
		if err != nil {
			fmt.Printf("Error getting default output device: %v\n", err)
			os.Exit(1)
		}
		mixerConfig.OutputDevice = dev
	}

	// Create mixer
	mixer, err := audio.NewMixer(mixerConfig)
	if err != nil {
		fmt.Printf("Error creating mixer: %v\n", err)
		os.Exit(1)
	}

	// Start mixer
	fmt.Println("\n=== Starting Audio Mixer ===")
	if err := mixer.Start(); err != nil {
		fmt.Printf("Error starting mixer: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Mixer started successfully!")
	fmt.Println("\nPress Ctrl+C to stop")
	fmt.Println("\nReal-time Monitoring:")
	fmt.Println("---------------------")

	// Setup signal handler for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	// Monitoring goroutine
	stopMonitor := make(chan struct{})
	go func() {
		ticker := time.NewTicker(500 * time.Millisecond)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				input1Level := mixer.GetInput1Level()
				input2Level := mixer.GetInput2Level()
				outputLevel := mixer.GetOutputLevel()
				latency := mixer.GetLatency()

				// Convert to dB for display
				input1DB := levelToDB(input1Level)
				input2DB := levelToDB(input2Level)
				outputDB := levelToDB(outputLevel)

				fmt.Printf("\r[Input1: %6.1f dB %s] [Input2: %6.1f dB %s] [Output: %6.1f dB %s] [Latency: %v]",
					input1DB, getLevelBar(input1Level, 20),
					input2DB, getLevelBar(input2Level, 20),
					outputDB, getLevelBar(outputLevel, 20),
					latency.Round(time.Microsecond))

			case <-stopMonitor:
				return
			}
		}
	}()

	// Wait for interrupt signal
	<-sigCh

	fmt.Println("\n\nShutting down...")
	close(stopMonitor)

	// Stop mixer
	if err := mixer.Stop(); err != nil {
		fmt.Printf("Error stopping mixer: %v\n", err)
	}

	fmt.Println("Goodbye!")
}

// readInt reads an integer from stdin with a default value
func readInt(reader *bufio.Reader, defaultValue int) int {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(input)
	if err != nil {
		fmt.Printf("Invalid input, using default: %d\n", defaultValue)
		return defaultValue
	}

	return value
}

// readFloat32 reads a float32 from stdin with a default value
func readFloat32(reader *bufio.Reader, defaultValue float32) float32 {
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		return defaultValue
	}

	value, err := strconv.ParseFloat(input, 32)
	if err != nil {
		fmt.Printf("Invalid input, using default: %.2f\n", defaultValue)
		return defaultValue
	}

	if value < 0 || value > 2.0 {
		fmt.Printf("Value out of range (0.0-2.0), using default: %.2f\n", defaultValue)
		return defaultValue
	}

	return float32(value)
}

// levelToDB converts linear level to decibels
func levelToDB(level float32) float32 {
	if level < 0.00001 {
		return -100.0
	}
	return 20.0 * float32(math.Log10(float64(level)))
}

// getLevelBar creates a visual level bar
func getLevelBar(level float32, width int) string {
	filled := int(level * float32(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	return bar
}
