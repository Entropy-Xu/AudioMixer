package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"github.com/entropy/audio-mixer/internal/audio"
	"github.com/entropy/audio-mixer/internal/config"
)

// App represents the GUI application
type App struct {
	fyneApp       fyne.App
	window        fyne.Window
	deviceManager *audio.DeviceManager
	mixer         *audio.Mixer
	configManager *config.ConfigManager
	cfg           *config.Config

	// UI elements
	input1Select  *widget.Select
	input2Select  *widget.Select
	outputSelect  *widget.Select
	input1Slider  *widget.Slider
	input2Slider  *widget.Slider
	masterSlider  *widget.Slider
	input1Label   *widget.Label
	input2Label   *widget.Label
	masterLabel   *widget.Label
	statusLabel   *widget.Label
	startButton   *widget.Button
	stopButton    *widget.Button
	input1Meter   *widget.ProgressBar
	input2Meter   *widget.ProgressBar
	outputMeter   *widget.ProgressBar
	latencyLabel  *widget.Label

	// State
	isRunning bool
}

// NewApp creates a new GUI application
func NewApp() (*App, error) {
	a := &App{
		fyneApp: app.New(),
	}

	// Initialize device manager
	a.deviceManager = audio.NewDeviceManager()
	if err := a.deviceManager.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize device manager: %w", err)
	}

	// Initialize config manager
	configManager, err := config.NewConfigManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize config manager: %w", err)
	}
	a.configManager = configManager

	// Load config
	cfg, err := configManager.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	a.cfg = cfg

	return a, nil
}

// Run starts the GUI application
func (a *App) Run() {
	a.window = a.fyneApp.NewWindow("Audio Mixer")
	a.window.Resize(fyne.NewSize(600, 700))

	// Build UI
	content := a.buildUI()
	a.window.SetContent(content)

	// Set close handler
	a.window.SetOnClosed(func() {
		a.cleanup()
	})

	a.window.ShowAndRun()
}

// buildUI creates the main UI layout
func (a *App) buildUI() fyne.CanvasObject {
	// Device selection section
	deviceSection := a.buildDeviceSection()

	// Volume control section
	volumeSection := a.buildVolumeSection()

	// Meters section
	metersSection := a.buildMetersSection()

	// Control buttons
	controlSection := a.buildControlSection()

	// Status bar
	a.statusLabel = widget.NewLabel("Ready")

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Audio Mixer Control Panel"),
		widget.NewSeparator(),
		deviceSection,
		widget.NewSeparator(),
		volumeSection,
		widget.NewSeparator(),
		metersSection,
		widget.NewSeparator(),
		controlSection,
		widget.NewSeparator(),
		a.statusLabel,
	)

	return container.NewPadded(content)
}

// buildDeviceSection creates device selection UI
func (a *App) buildDeviceSection() fyne.CanvasObject {
	// Get devices
	inputDevices, _ := a.deviceManager.GetInputDevices()
	outputDevices, _ := a.deviceManager.GetOutputDevices()

	// Build device names
	inputNames := make([]string, len(inputDevices))
	for i, dev := range inputDevices {
		inputNames[i] = fmt.Sprintf("[%d] %s", dev.Index, dev.Name)
	}

	outputNames := make([]string, len(outputDevices))
	for i, dev := range outputDevices {
		outputNames[i] = fmt.Sprintf("[%d] %s", dev.Index, dev.Name)
	}

	// Helper function to find device name by index
	findDeviceName := func(deviceIndex int, devices []*audio.DeviceInfo, names []string) string {
		for i, dev := range devices {
			if dev.Index == deviceIndex {
				return names[i]
			}
		}
		return ""
	}

	// Input 1 select
	a.input1Select = widget.NewSelect(inputNames, nil)
	a.input1Select.OnChanged = func(value string) {
		a.updateConfig()
	}
	if a.cfg.Input1DeviceIndex >= 0 {
		selectedName := findDeviceName(a.cfg.Input1DeviceIndex, inputDevices, inputNames)
		if selectedName != "" {
			a.input1Select.SetSelected(selectedName)
		}
	}

	// Input 2 select
	a.input2Select = widget.NewSelect(append([]string{"<None>"}, inputNames...), nil)
	a.input2Select.OnChanged = func(value string) {
		a.updateConfig()
	}
	if a.cfg.Input2DeviceIndex >= 0 {
		selectedName := findDeviceName(a.cfg.Input2DeviceIndex, inputDevices, inputNames)
		if selectedName != "" {
			a.input2Select.SetSelected(selectedName)
		}
	} else {
		a.input2Select.SetSelected("<None>")
	}

	// Output select
	a.outputSelect = widget.NewSelect(outputNames, nil)
	a.outputSelect.OnChanged = func(value string) {
		a.updateConfig()
	}
	if a.cfg.OutputDeviceIndex >= 0 {
		selectedName := findDeviceName(a.cfg.OutputDeviceIndex, outputDevices, outputNames)
		if selectedName != "" {
			a.outputSelect.SetSelected(selectedName)
		}
	}

	return container.NewVBox(
		widget.NewLabel("Device Selection"),
		container.New(layout.NewFormLayout(),
			widget.NewLabel("Input 1 (Microphone):"), a.input1Select,
			widget.NewLabel("Input 2 (App Audio):"), a.input2Select,
			widget.NewLabel("Output:"), a.outputSelect,
		),
	)
}

// buildVolumeSection creates volume control UI
func (a *App) buildVolumeSection() fyne.CanvasObject {
	// Input 1 gain
	a.input1Label = widget.NewLabel(fmt.Sprintf("Input 1 Gain: %.2f", a.cfg.Input1Gain))
	a.input1Slider = widget.NewSlider(0, 2.0)
	a.input1Slider.Value = float64(a.cfg.Input1Gain)
	a.input1Slider.Step = 0.01
	a.input1Slider.OnChanged = func(value float64) {
		a.input1Label.SetText(fmt.Sprintf("Input 1 Gain: %.2f", value))
		a.cfg.Input1Gain = float32(value)
		if a.isRunning && a.mixer != nil {
			a.mixer.SetInput1Gain(float32(value))
		}
	}

	// Input 2 gain
	a.input2Label = widget.NewLabel(fmt.Sprintf("Input 2 Gain: %.2f", a.cfg.Input2Gain))
	a.input2Slider = widget.NewSlider(0, 2.0)
	a.input2Slider.Value = float64(a.cfg.Input2Gain)
	a.input2Slider.Step = 0.01
	a.input2Slider.OnChanged = func(value float64) {
		a.input2Label.SetText(fmt.Sprintf("Input 2 Gain: %.2f", value))
		a.cfg.Input2Gain = float32(value)
		if a.isRunning && a.mixer != nil {
			a.mixer.SetInput2Gain(float32(value))
		}
	}

	// Master gain
	a.masterLabel = widget.NewLabel(fmt.Sprintf("Master Gain: %.2f", a.cfg.MasterGain))
	a.masterSlider = widget.NewSlider(0, 2.0)
	a.masterSlider.Value = float64(a.cfg.MasterGain)
	a.masterSlider.Step = 0.01
	a.masterSlider.OnChanged = func(value float64) {
		a.masterLabel.SetText(fmt.Sprintf("Master Gain: %.2f", value))
		a.cfg.MasterGain = float32(value)
		if a.isRunning && a.mixer != nil {
			a.mixer.SetMasterGain(float32(value))
		}
	}

	return container.NewVBox(
		widget.NewLabel("Volume Controls"),
		a.input1Label,
		a.input1Slider,
		a.input2Label,
		a.input2Slider,
		a.masterLabel,
		a.masterSlider,
	)
}

// buildMetersSection creates level meters UI
func (a *App) buildMetersSection() fyne.CanvasObject {
	a.input1Meter = widget.NewProgressBar()
	a.input2Meter = widget.NewProgressBar()
	a.outputMeter = widget.NewProgressBar()
	a.latencyLabel = widget.NewLabel("Latency: 0ms")

	return container.NewVBox(
		widget.NewLabel("Audio Levels"),
		widget.NewLabel("Input 1:"),
		a.input1Meter,
		widget.NewLabel("Input 2:"),
		a.input2Meter,
		widget.NewLabel("Output:"),
		a.outputMeter,
		a.latencyLabel,
	)
}

// buildControlSection creates start/stop buttons
func (a *App) buildControlSection() fyne.CanvasObject {
	a.startButton = widget.NewButton("Start Mixer", func() {
		a.startMixer()
	})

	a.stopButton = widget.NewButton("Stop Mixer", func() {
		a.stopMixer()
	})
	a.stopButton.Disable()

	return container.NewHBox(
		a.startButton,
		a.stopButton,
	)
}

// startMixer starts the audio mixer
func (a *App) startMixer() {
	if a.isRunning {
		return
	}

	// Save config
	a.updateConfig()
	if err := a.configManager.Save(a.cfg); err != nil {
		a.statusLabel.SetText(fmt.Sprintf("Error saving config: %v", err))
		return
	}

	// Setup mixer configuration
	mixerConfig := audio.DefaultMixerConfig()
	mixerConfig.SampleRate = a.cfg.SampleRate
	mixerConfig.BufferSize = a.cfg.BufferSize
	mixerConfig.Channels = a.cfg.Channels
	mixerConfig.Input1Gain = a.cfg.Input1Gain
	mixerConfig.Input2Gain = a.cfg.Input2Gain
	mixerConfig.MasterGain = a.cfg.MasterGain

	// Get device info
	if a.cfg.Input1DeviceIndex >= 0 {
		dev, err := a.deviceManager.GetDeviceByIndex(a.cfg.Input1DeviceIndex)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		mixerConfig.Input1Device = dev
	}

	if a.cfg.Input2DeviceIndex >= 0 {
		dev, err := a.deviceManager.GetDeviceByIndex(a.cfg.Input2DeviceIndex)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		mixerConfig.Input2Device = dev
	}

	if a.cfg.OutputDeviceIndex >= 0 {
		dev, err := a.deviceManager.GetDeviceByIndex(a.cfg.OutputDeviceIndex)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error: %v", err))
			return
		}
		mixerConfig.OutputDevice = dev
	}

	// Create mixer
	mixer, err := audio.NewMixer(mixerConfig)
	if err != nil {
		a.statusLabel.SetText(fmt.Sprintf("Error creating mixer: %v", err))
		return
	}
	a.mixer = mixer

	// Start mixer
	if err := a.mixer.Start(); err != nil {
		a.statusLabel.SetText(fmt.Sprintf("Error starting mixer: %v", err))
		return
	}

	a.isRunning = true
	a.startButton.Disable()
	a.stopButton.Enable()
	a.statusLabel.SetText("Mixer running")

	// Start meter update loop
	go a.updateMeters()
}

// stopMixer stops the audio mixer
func (a *App) stopMixer() {
	if !a.isRunning {
		return
	}

	if a.mixer != nil {
		if err := a.mixer.Stop(); err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error stopping mixer: %v", err))
		}
		a.mixer = nil
	}

	a.isRunning = false
	a.startButton.Enable()
	a.stopButton.Disable()
	a.statusLabel.SetText("Mixer stopped")

	// Reset meters
	a.input1Meter.SetValue(0)
	a.input2Meter.SetValue(0)
	a.outputMeter.SetValue(0)
	a.latencyLabel.SetText("Latency: 0ms")
}

// updateMeters updates the level meters
func (a *App) updateMeters() {
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for a.isRunning {
		<-ticker.C
		if a.mixer != nil {
			input1Level := a.mixer.GetInput1Level()
			input2Level := a.mixer.GetInput2Level()
			outputLevel := a.mixer.GetOutputLevel()
			latency := a.mixer.GetLatency()

			// Clamp to 0-1 range for display
			if input1Level > 1.0 {
				input1Level = 1.0
			}
			if input2Level > 1.0 {
				input2Level = 1.0
			}
			if outputLevel > 1.0 {
				outputLevel = 1.0
			}

			a.input1Meter.SetValue(float64(input1Level))
			a.input2Meter.SetValue(float64(input2Level))
			a.outputMeter.SetValue(float64(outputLevel))
			a.latencyLabel.SetText(fmt.Sprintf("Latency: %v", latency.Round(time.Microsecond)))
		}
	}
}

// updateConfig updates the config from UI selections
func (a *App) updateConfig() {
	// Parse device indices from selection
	if a.input1Select.Selected != "" && a.input1Select.Selected != "<None>" {
		var idx int
		fmt.Sscanf(a.input1Select.Selected, "[%d]", &idx)
		a.cfg.Input1DeviceIndex = idx
	}

	if a.input2Select.Selected != "" && a.input2Select.Selected != "<None>" {
		var idx int
		fmt.Sscanf(a.input2Select.Selected, "[%d]", &idx)
		a.cfg.Input2DeviceIndex = idx
	} else {
		a.cfg.Input2DeviceIndex = -2
	}

	if a.outputSelect.Selected != "" {
		var idx int
		fmt.Sscanf(a.outputSelect.Selected, "[%d]", &idx)
		a.cfg.OutputDeviceIndex = idx
	}
}

// cleanup performs cleanup when closing the app
func (a *App) cleanup() {
	if a.isRunning {
		a.stopMixer()
	}

	if a.deviceManager != nil {
		a.deviceManager.Terminate()
	}
}
