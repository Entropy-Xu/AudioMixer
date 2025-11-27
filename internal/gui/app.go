package gui

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"

	"github.com/entropy/audio-mixer/internal/audio"
	"github.com/entropy/audio-mixer/internal/config"
)

// App represents the GUI application
type App struct {
	fyneApp        fyne.App
	window         fyne.Window
	deviceManager  *audio.DeviceManager
	appCaptureManager *audio.ApplicationCaptureManager
	mixer          *audio.Mixer
	configManager  *config.ConfigManager
	cfg            *config.Config

	// UI elements
	input1Select        *widget.Select
	input2Select        *widget.Select
	appSelect           *widget.Select // Application audio selector
	outputNameEntry     *widget.Entry  // Custom output device name entry
	input1Slider        *widget.Slider
	input2Slider      *widget.Slider
	masterSlider      *widget.Slider
	input1Label       *widget.Label
	input2Label       *widget.Label
	masterLabel       *widget.Label
	statusLabel       *widget.Label
	startButton       *widget.Button
	stopButton        *widget.Button
	input1Meter       *widget.ProgressBar
	input2Meter       *widget.ProgressBar
	outputMeter       *widget.ProgressBar
	latencyLabel      *widget.Label
	fontSelect        *widget.Select
	fontStatus        *widget.Label

	// State
	isRunning bool
}

// NewApp creates a new GUI application
func NewApp() (*App, error) {
	a := &App{
		fyneApp: app.New(),
	}

	// Set custom theme with system font for better Unicode support
	a.fyneApp.Settings().SetTheme(&customTheme{})

	// Initialize device manager
	a.deviceManager = audio.NewDeviceManager()
	if err := a.deviceManager.Initialize(); err != nil {
		return nil, fmt.Errorf("failed to initialize device manager: %w", err)
	}

	// Initialize application capture manager
	appCapture, err := audio.NewApplicationCaptureManager()
	if err != nil {
		// 非致命错误，继续运行
		fmt.Fprintf(os.Stderr, "Warning: Application capture not available: %v\n", err)
	}
	a.appCaptureManager = appCapture

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

	// Font settings section
	fontSection := a.buildFontSection()

	// Status bar
	a.statusLabel = widget.NewLabel("Ready")

	// Main layout
	content := container.NewVBox(
		widget.NewLabel("Audio Mixer"),
		widget.NewSeparator(),
		deviceSection,
		widget.NewSeparator(),
		volumeSection,
		widget.NewSeparator(),
		metersSection,
		widget.NewSeparator(),
		controlSection,
		widget.NewSeparator(),
		fontSection,
		widget.NewSeparator(),
		a.statusLabel,
	)

	return container.NewPadded(content)
}

// buildDeviceSection creates device selection UI
func (a *App) buildDeviceSection() fyne.CanvasObject {
	// Get devices
	inputDevices, _ := a.deviceManager.GetInputDevices()

	// Get loopback devices for Input2 and Output
	loopbackDevices, _ := audio.ListLoopbackDevices(a.deviceManager)

	// Build Input 1 device names (microphone/line input)
	inputNames := make([]string, len(inputDevices))
	for i, dev := range inputDevices {
		inputNames[i] = fmt.Sprintf("[%d] %s", dev.Index, dev.Name)
	}

	// Build Input 2 options (loopback devices for system audio)
	input2Names := []string{"<Auto Detect Loopback>"}
	var input2Devices []*audio.DeviceInfo
	for _, lb := range loopbackDevices {
		// Convert loopback device to DeviceInfo for consistency
		devInfo := &audio.DeviceInfo{
			Name: lb.Name,
		}
		// Find the actual device index
		allDevs, _ := a.deviceManager.ListDevices()
		for _, d := range allDevs {
			if d.Name == lb.Name {
				devInfo.Index = d.Index
				break
			}
		}
		input2Devices = append(input2Devices, devInfo)
		input2Names = append(input2Names, fmt.Sprintf("[Virtual] %s", lb.Name))
	}

	// Build Output options (loopback/virtual devices only)
	outputNames := []string{"<Auto Detect Virtual Output>"}
	var outputDevices []*audio.DeviceInfo
	for _, lb := range loopbackDevices {
		devInfo := &audio.DeviceInfo{
			Name: lb.Name,
		}
		allDevs, _ := a.deviceManager.ListDevices()
		for _, d := range allDevs {
			if d.Name == lb.Name {
				devInfo.Index = d.Index
				break
			}
		}
		outputDevices = append(outputDevices, devInfo)
		outputNames = append(outputNames, fmt.Sprintf("[Virtual] %s", lb.Name))
	}

	// Helper function to find device name by index
	findDeviceName := func(deviceIndex int, devices []*audio.DeviceInfo, names []string) string {
		for i, dev := range devices {
			if dev.Index == deviceIndex {
				if i < len(names) {
					return names[i]
				}
			}
		}
		return ""
	}

	// Input 1 select (Microphone/Line Input)
	a.input1Select = widget.NewSelect(inputNames, nil)
	if a.cfg.Input1DeviceIndex >= 0 {
		selectedName := findDeviceName(a.cfg.Input1DeviceIndex, inputDevices, inputNames)
		if selectedName != "" {
			a.input1Select.SetSelected(selectedName)
		}
	}
	a.input1Select.OnChanged = func(value string) {
		if value != "" {
			a.updateConfig()
		}
	}

	// Input 2 select (System Audio via Loopback)
	a.input2Select = widget.NewSelect(input2Names, nil)
	if a.cfg.Input2DeviceIndex >= 0 {
		selectedName := findDeviceName(a.cfg.Input2DeviceIndex, input2Devices, input2Names[1:])
		if selectedName != "" {
			a.input2Select.SetSelected(selectedName)
		} else {
			a.input2Select.SetSelected("<Auto Detect Loopback>")
		}
	} else {
		a.input2Select.SetSelected("<Auto Detect Loopback>")
	}
	a.input2Select.OnChanged = func(value string) {
		if value != "" {
			a.updateConfig()
		}
	}

	// Output select (Virtual Output Device) - removed, using custom name instead

	// Custom output device name entry
	a.outputNameEntry = widget.NewEntry()
	a.outputNameEntry.SetPlaceHolder("输入虚拟设备名称，如: BlackHole 2ch")
	if a.cfg.LoopbackDeviceName != "" {
		a.outputNameEntry.SetText(a.cfg.LoopbackDeviceName)
	} else {
		a.outputNameEntry.SetText("BlackHole")
	}
	a.outputNameEntry.OnChanged = func(value string) {
		a.cfg.LoopbackDeviceName = value
		a.updateConfig()
	}

	// Detect button to find the device
	detectButton := widget.NewButton("检测设备 (Detect)", func() {
		if a.outputNameEntry.Text == "" {
			a.statusLabel.SetText("请先输入设备名称")
			return
		}

		// Try to find device by name
		loopback, err := audio.GetLoopbackDeviceByName(a.deviceManager, a.outputNameEntry.Text)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("未找到设备: %v", err))
		} else {
			a.statusLabel.SetText(fmt.Sprintf("✓ 找到设备: %s", loopback.Name))
		}
	})

	// Application selector (for capturing specific app audio)
	var appSection *fyne.Container
	if a.appCaptureManager != nil {
		appNames := []string{"<使用虚拟设备方案>"}

		// Get list of applications
		apps, err := a.appCaptureManager.ListApplications()
		if err == nil {
			for _, app := range apps {
				appNames = append(appNames, audio.GetFriendlyName(app.DisplayName))
			}
		}

		a.appSelect = widget.NewSelect(appNames, func(selected string) {
			// TODO: Implement application capture
			if selected != "<使用虚拟设备方案>" {
				a.statusLabel.SetText(fmt.Sprintf("注意: 应用音频捕获功能仍在开发中"))
			}
		})
		a.appSelect.SetSelected("<使用虚拟设备方案>")

		refreshButton := widget.NewButton("🔄 刷新", func() {
			// Refresh application list
			apps, err := a.appCaptureManager.ListApplications()
			if err != nil {
				a.statusLabel.SetText(fmt.Sprintf("刷新失败: %v", err))
				return
			}

			newNames := []string{"<使用虚拟设备方案>"}
			for _, app := range apps {
				newNames = append(newNames, audio.GetFriendlyName(app.DisplayName))
			}
			a.appSelect.Options = newNames
			a.appSelect.Refresh()
			a.statusLabel.SetText(fmt.Sprintf("已刷新，找到 %d 个应用", len(apps)))
		})

		appSection = container.NewVBox(
			widget.NewLabel("或者，捕获特定应用音频 (实验性):"),
			container.NewBorder(nil, nil, nil, refreshButton, a.appSelect),
		)
	}

	// Info label explaining the setup
	infoText := "提示:\n" +
		"• Input 1: 麦克风/线路输入\n" +
		"• Input 2: 系统音频 (需要虚拟设备,如 BlackHole/VB-Cable)\n" +
		"• Output: 输入虚拟设备名称 (混音后输出到该设备)\n" +
		"• 常见设备名: BlackHole 2ch, CABLE-B Input, VB-Cable"
	infoLabel := widget.NewLabel(infoText)

	sections := []fyne.CanvasObject{
		widget.NewLabel("设备配置 (Devices)"),
		container.New(layout.NewFormLayout(),
			widget.NewLabel("Input 1 (麦克风):"), a.input1Select,
			widget.NewLabel("Input 2 (系统音频):"), a.input2Select,
		),
		widget.NewLabel("Output (虚拟输出设备名称):"),
		container.NewBorder(nil, nil, nil, detectButton, a.outputNameEntry),
	}

	if appSection != nil {
		sections = append(sections, appSection)
	}

	sections = append(sections, infoLabel)

	return container.NewVBox(sections...)
}

// buildVolumeSection creates volume control UI
func (a *App) buildVolumeSection() fyne.CanvasObject {
	// Input 1 gain
	a.input1Label = widget.NewLabel(fmt.Sprintf("Input 1: %.2f", a.cfg.Input1Gain))
	a.input1Slider = widget.NewSlider(0, 2.0)
	a.input1Slider.Value = float64(a.cfg.Input1Gain)
	a.input1Slider.Step = 0.01
	a.input1Slider.OnChanged = func(value float64) {
		a.input1Label.SetText(fmt.Sprintf("Input 1: %.2f", value))
		a.cfg.Input1Gain = float32(value)
		if a.isRunning && a.mixer != nil {
			a.mixer.SetInput1Gain(float32(value))
		}
	}

	// Input 2 gain
	a.input2Label = widget.NewLabel(fmt.Sprintf("Input 2: %.2f", a.cfg.Input2Gain))
	a.input2Slider = widget.NewSlider(0, 2.0)
	a.input2Slider.Value = float64(a.cfg.Input2Gain)
	a.input2Slider.Step = 0.01
	a.input2Slider.OnChanged = func(value float64) {
		a.input2Label.SetText(fmt.Sprintf("Input 2: %.2f", value))
		a.cfg.Input2Gain = float32(value)
		if a.isRunning && a.mixer != nil {
			a.mixer.SetInput2Gain(float32(value))
		}
	}

	// Master gain
	a.masterLabel = widget.NewLabel(fmt.Sprintf("Master: %.2f", a.cfg.MasterGain))
	a.masterSlider = widget.NewSlider(0, 2.0)
	a.masterSlider.Value = float64(a.cfg.MasterGain)
	a.masterSlider.Step = 0.01
	a.masterSlider.OnChanged = func(value float64) {
		a.masterLabel.SetText(fmt.Sprintf("Master: %.2f", value))
		a.cfg.MasterGain = float32(value)
		if a.isRunning && a.mixer != nil {
			a.mixer.SetMasterGain(float32(value))
		}
	}

	return container.NewVBox(
		widget.NewLabel("Volume (0.00-2.00)"),
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
		widget.NewLabel("Levels"),
		widget.NewLabel("In1:"),
		a.input1Meter,
		widget.NewLabel("In2:"),
		a.input2Meter,
		widget.NewLabel("Out:"),
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

// buildFontSection creates font selection UI
func (a *App) buildFontSection() fyne.CanvasObject {
	// Get available fonts
	fontPaths := GetDefaultFontPaths()
	fontOptions := []string{"Default (Built-in)"}
	fontMap := make(map[string]string)

	fontMap["Default (Built-in)"] = ""

	// Check which fonts are available
	for _, path := range fontPaths {
		if _, err := os.Stat(path); err == nil {
			// Extract friendly name from path
			name := getFriendlyFontName(path)
			fontOptions = append(fontOptions, name)
			fontMap[name] = path
		}
	}

	// Create font select dropdown
	a.fontSelect = widget.NewSelect(fontOptions, nil)

	// Set current font
	currentFont := os.Getenv("FYNE_FONT")
	selectedOption := "Default (Built-in)"
	for name, path := range fontMap {
		if path == currentFont {
			selectedOption = name
			break
		}
	}
	a.fontSelect.SetSelected(selectedOption)

	// Font status label
	a.fontStatus = widget.NewLabel("")
	if currentFont != "" {
		a.fontStatus.SetText(fmt.Sprintf("Current: %s", currentFont))
	}

	// Set callback after initial value
	a.fontSelect.OnChanged = func(selected string) {
		if selected == "" {
			return
		}

		fontPath := fontMap[selected]
		if err := LoadCustomFont(fontPath); err != nil {
			a.fontStatus.SetText(fmt.Sprintf("Error: %v", err))
			return
		}

		if fontPath == "" {
			a.fontStatus.SetText("Using default font")
		} else {
			a.fontStatus.SetText(fmt.Sprintf("Loaded: %s", fontPath))
		}

		// Show restart message
		a.statusLabel.SetText("Font changed. Please restart the app for changes to take effect.")
	}

	// Add custom font button
	customButton := widget.NewButton("Custom Font...", func() {
		// This would open a file dialog, but for simplicity we'll show a message
		a.statusLabel.SetText("Use: ./audio-mixer-gui -font /path/to/font.ttf")
	})

	return container.NewVBox(
		widget.NewLabel("Font Settings (for CJK characters)"),
		container.New(layout.NewFormLayout(),
			widget.NewLabel("Font:"), a.fontSelect,
		),
		customButton,
		a.fontStatus,
	)
}

// getFriendlyFontName extracts a friendly name from font path
func getFriendlyFontName(path string) string {
	switch path {
	case "/System/Library/Fonts/Supplemental/Arial Unicode.ttf":
		return "Arial Unicode MS (Recommended)"
	case "/Library/Fonts/Arial Unicode.ttf":
		return "Arial Unicode MS"
	case "/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttf":
		return "Noto Sans CJK"
	case "/usr/share/fonts/opentype/noto/NotoSansCJKsc-Regular.otf":
		return "Noto Sans CJK SC (OpenType)"
	case "C:\\Windows\\Fonts\\msyh.ttf":
		return "Microsoft YaHei"
	case "C:\\Windows\\Fonts\\simsun.ttf":
		return "SimSun"
	case "C:\\Windows\\Fonts\\arial.ttf":
		return "Arial"
	default:
		// Check if it's a user font (likely custom installed CJK font)
		if strings.Contains(path, "PingFangSC") {
			return "PingFang SC (User Font)"
		}
		if strings.Contains(path, "SourceHanSans") {
			return "Source Han Sans (User Font)"
		}
		if strings.Contains(path, "NotoSansCJK") {
			return "Noto Sans CJK (User Font)"
		}

		// Extract filename
		parts := strings.Split(path, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
		return path
	}
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
	mixerConfig.UseVirtualOutput = a.cfg.UseVirtualOutput

	// Get Input 1 device (microphone/line input)
	if a.cfg.Input1DeviceIndex >= 0 {
		dev, err := a.deviceManager.GetDeviceByIndex(a.cfg.Input1DeviceIndex)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error getting Input 1: %v", err))
			return
		}
		mixerConfig.Input1Device = dev
	} else {
		// Use default input device
		dev, err := a.deviceManager.GetDefaultInputDevice()
		if err == nil {
			mixerConfig.Input1Device = dev
		}
	}

	// Get Input 2 device (system audio via loopback)
	if a.cfg.Input2DeviceIndex >= 0 {
		dev, err := a.deviceManager.GetDeviceByIndex(a.cfg.Input2DeviceIndex)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("Error getting Input 2: %v", err))
			return
		}
		mixerConfig.Input2Device = dev
	} else {
		// Auto-detect loopback device for system audio
		loopback, err := audio.FindLoopbackDevice(a.deviceManager)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("未找到虚拟设备! 请安装 BlackHole: %v", err))
			return
		}
		mixerConfig.Input2Device = loopback.Device
		a.statusLabel.SetText(fmt.Sprintf("自动检测到虚拟设备: %s", loopback.Name))
	}

	// Get Output device (virtual device by custom name)
	if a.cfg.LoopbackDeviceName != "" {
		// Use custom device name
		loopback, err := audio.GetLoopbackDeviceByName(a.deviceManager, a.cfg.LoopbackDeviceName)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("未找到设备 '%s': %v\n请检查设备名称或点击'检测设备'", a.cfg.LoopbackDeviceName, err))
			return
		}
		mixerConfig.OutputDevice = loopback.Device
		a.statusLabel.SetText(fmt.Sprintf("使用虚拟输出: %s", loopback.Name))
	} else {
		// Auto-detect virtual output device (fallback)
		loopback, err := audio.FindLoopbackDevice(a.deviceManager)
		if err != nil {
			a.statusLabel.SetText(fmt.Sprintf("未找到虚拟输出设备! 请安装 BlackHole 或输入设备名称: %v", err))
			return
		}
		mixerConfig.OutputDevice = loopback.Device
		a.statusLabel.SetText(fmt.Sprintf("自动检测到虚拟输出: %s", loopback.Name))
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
	a.statusLabel.SetText("混音器运行中 (Mixer running)")

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
	if a.input1Select != nil && a.input1Select.Selected != "" && a.input1Select.Selected != "<None>" {
		var idx int
		fmt.Sscanf(a.input1Select.Selected, "[%d]", &idx)
		a.cfg.Input1DeviceIndex = idx
	}

	if a.input2Select != nil {
		if a.input2Select.Selected != "" && a.input2Select.Selected != "<None>" {
			var idx int
			fmt.Sscanf(a.input2Select.Selected, "[%d]", &idx)
			a.cfg.Input2DeviceIndex = idx
		} else {
			a.cfg.Input2DeviceIndex = -2
		}
	}

	// Output device is now configured via custom name (outputNameEntry)
	// LoopbackDeviceName is already updated in outputNameEntry.OnChanged
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

// customTheme provides a theme with better font fallback for Unicode characters
type customTheme struct{}

func (c *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.DefaultTheme().Color(name, variant)
}

func (c *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (c *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	// Try to use system font on macOS for better Unicode support
	if style.Monospace {
		return theme.DefaultTheme().Font(style)
	}

	// On macOS, check for PingFang or Heiti font for Chinese support
	if _, err := os.Stat("/System/Library/Fonts/PingFang.ttc"); err == nil {
		// System has PingFang font, but we can't load it directly in Fyne
		// Fall back to default theme which should handle system fonts
		return theme.DefaultTheme().Font(style)
	}

	return theme.DefaultTheme().Font(style)
}

func (c *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
