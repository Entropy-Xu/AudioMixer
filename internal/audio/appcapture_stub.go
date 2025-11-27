// +build !windows

package audio

import (
	"fmt"
	"runtime"
)

// stubCaptureImpl 非 Windows 平台的存根实现
type stubCaptureImpl struct{}

// newPlatformCaptureImpl 创建平台特定实现
func newPlatformCaptureImpl() (applicationCaptureImpl, error) {
	return &stubCaptureImpl{}, nil
}

// ListApplications 列出所有正在播放音频的应用
func (s *stubCaptureImpl) ListApplications() ([]*ApplicationInfo, error) {
	// macOS 和 Linux 平台的应用音频捕获需要不同的实现
	//
	// macOS: 使用 ScreenCaptureKit API (需要 macOS 13+)
	// Linux: 使用 PulseAudio API
	//
	// 目前返回占位信息

	return []*ApplicationInfo{
		{
			ProcessID:   0,
			ProcessName: fmt.Sprintf("%s_stub", runtime.GOOS),
			DisplayName: fmt.Sprintf("⚠️ 应用音频捕获功能暂不支持 %s", runtime.GOOS),
			IsPlaying:   false,
			Volume:      0,
		},
		{
			ProcessID:   1,
			ProcessName: "workaround",
			DisplayName: "💡 请使用虚拟设备方案（参见 QUICK_SETUP_GUIDE.md）",
			IsPlaying:   false,
			Volume:      0,
		},
	}, nil
}

// StartCapture 开始捕获指定应用的音频
func (s *stubCaptureImpl) StartCapture(processID uint32, callback func([]float32)) error {
	return fmt.Errorf("application audio capture is not supported on %s", runtime.GOOS)
}

// StopCapture 停止捕获
func (s *stubCaptureImpl) StopCapture() error {
	return nil
}

// IsCapturing 是否正在捕获
func (s *stubCaptureImpl) IsCapturing() bool {
	return false
}
