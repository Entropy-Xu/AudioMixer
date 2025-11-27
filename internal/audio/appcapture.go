package audio

import (
	"fmt"
)

// ApplicationInfo 应用程序音频信息
type ApplicationInfo struct {
	ProcessID   uint32  // 进程ID
	ProcessName string  // 进程名称（如 "spotify.exe"）
	DisplayName string  // 显示名称（如 "Spotify"）
	IconPath    string  // 图标路径
	IsPlaying   bool    // 是否正在播放音频
	Volume      float32 // 当前音量 (0.0 - 1.0)
}

// ApplicationCaptureManager 应用音频捕获管理器
type ApplicationCaptureManager struct {
	// 平台特定实现
	impl applicationCaptureImpl
}

// applicationCaptureImpl 平台特定接口
type applicationCaptureImpl interface {
	// ListApplications 列出所有正在播放音频的应用
	ListApplications() ([]*ApplicationInfo, error)

	// StartCapture 开始捕获指定应用的音频
	StartCapture(processID uint32, callback func([]float32)) error

	// StopCapture 停止捕获
	StopCapture() error

	// IsCapturing 是否正在捕获
	IsCapturing() bool
}

// NewApplicationCaptureManager 创建应用捕获管理器
func NewApplicationCaptureManager() (*ApplicationCaptureManager, error) {
	impl, err := newPlatformCaptureImpl()
	if err != nil {
		return nil, err
	}

	return &ApplicationCaptureManager{
		impl: impl,
	}, nil
}

// ListApplications 列出所有正在播放音频的应用
func (m *ApplicationCaptureManager) ListApplications() ([]*ApplicationInfo, error) {
	return m.impl.ListApplications()
}

// StartCapture 开始捕获指定应用的音频
func (m *ApplicationCaptureManager) StartCapture(processID uint32, callback func([]float32)) error {
	return m.impl.StartCapture(processID, callback)
}

// StopCapture 停止捕获
func (m *ApplicationCaptureManager) StopCapture() error {
	return m.impl.StopCapture()
}

// IsCapturing 是否正在捕获
func (m *ApplicationCaptureManager) IsCapturing() bool {
	return m.impl.IsCapturing()
}

// GetFriendlyName 获取应用的友好名称
func GetFriendlyName(processName string) string {
	// 常见应用的友好名称映射
	friendlyNames := map[string]string{
		"spotify.exe":      "🎵 Spotify",
		"chrome.exe":       "🌐 Google Chrome",
		"firefox.exe":      "🦊 Firefox",
		"msedge.exe":       "🌊 Microsoft Edge",
		"discord.exe":      "💬 Discord",
		"vlc.exe":          "▶️ VLC Media Player",
		"wmplayer.exe":     "▶️ Windows Media Player",
		"itunes.exe":       "🎵 iTunes",
		"musicbee.exe":     "🎵 MusicBee",
		"foobar2000.exe":   "🎵 foobar2000",
		"steam.exe":        "🎮 Steam",
		"obs64.exe":        "📹 OBS Studio",
		"teams.exe":        "👥 Microsoft Teams",
		"slack.exe":        "💼 Slack",
		"zoom.exe":         "📹 Zoom",
		"skype.exe":        "📞 Skype",
	}

	if friendly, ok := friendlyNames[processName]; ok {
		return friendly
	}

	// 移除 .exe 后缀
	if len(processName) > 4 && processName[len(processName)-4:] == ".exe" {
		return processName[:len(processName)-4]
	}

	return processName
}
