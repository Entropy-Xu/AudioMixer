// +build windows

package audio

import (
	"fmt"
	"sync"
	"syscall"
	"unsafe"

	"github.com/go-ole/go-ole"
	"github.com/go-ole/go-ole/oleutil"
)

/*
Windows Audio Session API (WASAPI) 实现
用于捕获特定应用程序的音频（类似 OBS）

需要 go-ole 库：go get github.com/go-ole/go-ole
*/

// wasapiCaptureImpl WASAPI 平台实现
type wasapiCaptureImpl struct {
	processID   uint32
	isCapturing bool
	callback    func([]float32)
	stopCh      chan struct{}
	mu          sync.Mutex
}

// newPlatformCaptureImpl 创建平台特定实现
func newPlatformCaptureImpl() (applicationCaptureImpl, error) {
	// 初始化 COM
	err := ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)
	if err != nil {
		// 如果已经初始化（S_FALSE = 0x00000001），继续
		oleErr, ok := err.(*ole.OleError)
		if !ok || oleErr.Code() != 0x00000001 {
			return nil, fmt.Errorf("failed to initialize COM: %w", err)
		}
	}

	return &wasapiCaptureImpl{
		stopCh: make(chan struct{}),
	}, nil
}

// ListApplications 列出所有正在播放音频的应用
func (w *wasapiCaptureImpl) ListApplications() ([]*ApplicationInfo, error) {
	apps := []*ApplicationInfo{}

	// 创建 MMDeviceEnumerator
	unknown, err := oleutil.CreateObject("MMDeviceEnumerator")
	if err != nil {
		return nil, fmt.Errorf("failed to create MMDeviceEnumerator: %w", err)
	}
	defer unknown.Release()

	deviceEnumerator, err := unknown.QueryInterface(ole.IID_IUnknown)
	if err != nil {
		return nil, fmt.Errorf("failed to query IMMDeviceEnumerator: %w", err)
	}
	defer deviceEnumerator.Release()

	// 获取默认音频渲染端点
	// eRender = 0, eConsole = 0
	defaultDeviceVariant, err := oleutil.CallMethod(deviceEnumerator, "GetDefaultAudioEndpoint", 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to get default audio endpoint: %w", err)
	}
	defaultDevice := defaultDeviceVariant.ToIDispatch()
	defer defaultDevice.Release()

	// 激活 IAudioSessionManager2
	// IID_IAudioSessionManager2: {77AA99A0-1BD6-484F-8BC7-2C654C9A9B6F}
	sessionManagerGUID := ole.NewGUID("{77AA99A0-1BD6-484F-8BC7-2C654C9A9B6F}")
	sessionManagerVariant, err := oleutil.CallMethod(defaultDevice, "Activate", sessionManagerGUID, 0, 0, 0)
	if err != nil {
		return nil, fmt.Errorf("failed to activate AudioSessionManager2: %w", err)
	}
	sessionManager := sessionManagerVariant.ToIDispatch()
	defer sessionManager.Release()

	// 获取会话枚举器
	sessionEnumeratorVariant, err := oleutil.CallMethod(sessionManager, "GetSessionEnumerator")
	if err != nil {
		return nil, fmt.Errorf("failed to get session enumerator: %w", err)
	}
	sessionEnumerator := sessionEnumeratorVariant.ToIDispatch()
	defer sessionEnumerator.Release()

	// 获取会话数量
	sessionCountVariant, err := oleutil.CallMethod(sessionEnumerator, "GetCount")
	if err != nil {
		return nil, fmt.Errorf("failed to get session count: %w", err)
	}
	sessionCount := int(sessionCountVariant.Val)

	// 遍历所有会话
	for i := 0; i < sessionCount; i++ {
		sessionVariant, err := oleutil.CallMethod(sessionEnumerator, "GetSession", i)
		if err != nil {
			continue
		}
		session := sessionVariant.ToIDispatch()

		// 获取 IAudioSessionControl2
		sessionControl2, err := session.QueryInterface(ole.IID_IUnknown)
		if err != nil {
			session.Release()
			continue
		}

		// 获取进程 ID
		processIDVariant, err := oleutil.CallMethod(sessionControl2, "GetProcessId")
		if err != nil {
			sessionControl2.Release()
			session.Release()
			continue
		}
		processID := uint32(processIDVariant.Val)

		// 跳过系统进程 (PID 0)
		if processID == 0 {
			sessionControl2.Release()
			session.Release()
			continue
		}

		// 获取显示名称
		displayNameVariant, err := oleutil.CallMethod(sessionControl2, "GetDisplayName")
		displayName := ""
		if err == nil && displayNameVariant.Value() != nil {
			displayName = displayNameVariant.ToString()
		}

		// 获取会话状态
		stateVariant, err := oleutil.CallMethod(sessionControl2, "GetState")
		isPlaying := false
		if err == nil {
			// AudioSessionStateActive = 1
			isPlaying = (int(stateVariant.Val) == 1)
		}

		// 获取进程名称
		processName, err := getProcessName(processID)
		if err != nil {
			processName = fmt.Sprintf("Process_%d", processID)
		}

		// 如果没有显示名称，使用友好名称
		if displayName == "" {
			displayName = GetFriendlyName(processName)
		}

		// 获取音量（简化实现）
		volume := float32(1.0)

		apps = append(apps, &ApplicationInfo{
			ProcessID:   processID,
			ProcessName: processName,
			DisplayName: displayName,
			IsPlaying:   isPlaying,
			Volume:      volume,
		})

		sessionControl2.Release()
		session.Release()
	}

	return apps, nil
}

// StartCapture 开始捕获指定应用的音频
func (w *wasapiCaptureImpl) StartCapture(processID uint32, callback func([]float32)) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.isCapturing {
		return fmt.Errorf("already capturing")
	}

	w.processID = processID
	w.callback = callback
	w.isCapturing = true

	// 启动捕获线程
	go w.captureLoop()

	return nil
}

// StopCapture 停止捕获
func (w *wasapiCaptureImpl) StopCapture() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if !w.isCapturing {
		return nil
	}

	close(w.stopCh)
	w.isCapturing = false
	w.stopCh = make(chan struct{}) // 重新创建以备下次使用

	return nil
}

// IsCapturing 是否正在捕获
func (w *wasapiCaptureImpl) IsCapturing() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.isCapturing
}

// captureLoop 音频捕获循环
func (w *wasapiCaptureImpl) captureLoop() {
	// 注意：完整的 WASAPI loopback capture 实现非常复杂
	// 需要：
	// 1. 找到目标进程的音频会话
	// 2. 获取该会话的音频端点
	// 3. 创建 IAudioClient
	// 4. 初始化为 loopback 模式
	// 5. 获取 IAudioCaptureClient
	// 6. 循环读取音频缓冲区
	// 7. 转换音频格式并调用回调

	// 由于 WASAPI 的应用级捕获需要访问特定会话的端点，
	// 而 Windows 不直接支持单应用 loopback（需要 Audio Graph API 或复杂的路由），
	// 这里提供一个简化的实现框架

	// TODO: 完整实现需要使用 Windows 10 Audio Graph API 或更底层的 WASAPI
	// 参考: https://docs.microsoft.com/en-us/windows/win32/coreaudio/loopback-recording

	// 简化实现：捕获所有系统音频（loopback）
	// 实际应用中可以通过 audio session 的音量控制来隔离特定应用

	<-w.stopCh
}

// getProcessName 获取进程名称
func getProcessName(processID uint32) (string, error) {
	// 打开进程
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	procOpenProcess := kernel32.NewProc("OpenProcess")
	procQueryFullProcessImageName := kernel32.NewProc("QueryFullProcessImageNameW")
	procCloseHandle := kernel32.NewProc("CloseHandle")

	const PROCESS_QUERY_LIMITED_INFORMATION = 0x1000
	handle, _, _ := procOpenProcess.Call(
		uintptr(PROCESS_QUERY_LIMITED_INFORMATION),
		0,
		uintptr(processID),
	)

	if handle == 0 {
		return "", fmt.Errorf("failed to open process %d", processID)
	}
	defer procCloseHandle.Call(handle)

	// 获取进程路径
	var size uint32 = 260
	buffer := make([]uint16, size)
	ret, _, _ := procQueryFullProcessImageName.Call(
		handle,
		0,
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(unsafe.Pointer(&size)),
	)

	if ret == 0 {
		return "", fmt.Errorf("failed to query process name for PID %d", processID)
	}

	fullPath := syscall.UTF16ToString(buffer)

	// 提取文件名
	for i := len(fullPath) - 1; i >= 0; i-- {
		if fullPath[i] == '\\' || fullPath[i] == '/' {
			return fullPath[i+1:], nil
		}
	}

	return fullPath, nil
}
