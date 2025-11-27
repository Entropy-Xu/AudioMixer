# WASAPI 实现说明

## 最新更新 (2025-11-27)

### 已完成的功能

✅ **应用音频会话枚举**
- 完整实现了 Windows WASAPI 应用枚举功能
- 使用 `github.com/go-ole/go-ole` 库简化 COM 接口调用
- 可以列出所有正在播放音频的应用程序

### 实现细节

#### 1. COM 接口调用流程

```go
// 1. 初始化 COM
ole.CoInitializeEx(0, ole.COINIT_APARTMENTTHREADED)

// 2. 创建设备枚举器
unknown, _ := oleutil.CreateObject("MMDeviceEnumerator")
deviceEnumerator, _ := unknown.QueryInterface(ole.IID_IUnknown)

// 3. 获取默认音频渲染端点
defaultDevice := oleutil.CallMethod(deviceEnumerator, "GetDefaultAudioEndpoint", 0, 0)

// 4. 激活 IAudioSessionManager2
sessionManagerGUID := ole.NewGUID("{77AA99A0-1BD6-484F-8BC7-2C654C9A9B6F}")
sessionManager := oleutil.CallMethod(defaultDevice, "Activate", sessionManagerGUID, 0, 0, 0)

// 5. 获取会话枚举器
sessionEnumerator := oleutil.CallMethod(sessionManager, "GetSessionEnumerator")

// 6. 遍历所有音频会话
sessionCount := oleutil.CallMethod(sessionEnumerator, "GetCount")
for i := 0; i < sessionCount; i++ {
    session := oleutil.CallMethod(sessionEnumerator, "GetSession", i)
    sessionControl2, _ := session.QueryInterface(ole.IID_IUnknown)

    // 获取进程 ID
    processID := oleutil.CallMethod(sessionControl2, "GetProcessId")

    // 获取显示名称
    displayName := oleutil.CallMethod(sessionControl2, "GetDisplayName")

    // 获取会话状态
    state := oleutil.CallMethod(sessionControl2, "GetState")
}
```

#### 2. 进程信息获取

实现了 `getProcessName()` 函数，使用 Windows API 获取进程名称：

```go
func getProcessName(processID uint32) (string, error) {
    // 使用 kernel32.dll
    // OpenProcess - 打开进程句柄
    // QueryFullProcessImageNameW - 获取进程完整路径
    // 提取文件名
}
```

#### 3. 应用信息结构

```go
type ApplicationInfo struct {
    ProcessID   uint32  // 进程 ID
    ProcessName string  // 进程名（如 "spotify.exe"）
    DisplayName string  // 友好显示名（如 "🎵 Spotify"）
    IsPlaying   bool    // 是否正在播放音频
    Volume      float32 // 音量 (0.0 - 1.0)
}
```

### 使用方法

#### Windows 平台

1. **安装依赖**
   ```bash
   go get github.com/go-ole/go-ole
   ```

2. **编译**
   ```bash
   go build -o audio-mixer-gui.exe cmd/gui/main.go
   ```

3. **运行应用**
   - 启动 Audio Mixer GUI
   - 在"捕获特定应用音频"下拉框中点击"刷新"按钮
   - 下拉框会显示所有正在播放音频的应用

4. **查看应用列表**
   - ✅ 显示进程名（如 "spotify.exe"）
   - ✅ 显示友好名称（如 "🎵 Spotify"）
   - ✅ 显示播放状态（活动/非活动）
   - ✅ 实时更新（点击刷新）

#### macOS/Linux 平台

当前使用存根实现：
- 显示提示信息：不支持应用音频捕获
- 建议使用虚拟设备方案（BlackHole、PulseAudio）

### 技术限制说明

#### Windows 应用级 Loopback 捕获

**问题**: Windows WASAPI 不直接支持单应用的 loopback 音频捕获

**原因**:
1. WASAPI Loopback 模式只能捕获整个音频端点的输出
2. 不能直接捕获单个 audio session 的输出
3. 需要使用 Windows 10+ 的 Audio Graph API 或复杂的音频路由

**当前实现**:
- ✅ 可以枚举应用和会话
- ✅ 可以获取会话状态和音量
- ⚠️ 音频捕获框架已建立，但需要配合虚拟设备使用

**推荐方案**:
使用 VB-Cable 虚拟设备方案（详见 WINDOWS_SETUP_GUIDE.md）：

1. 安装 VB-Cable
2. 将目标应用的输出设为 CABLE Input
3. Audio Mixer 从 CABLE Output 读取
4. 混音后输出到 CABLE-B Input
5. OBS/Zoom 从 CABLE-B Output 读取

这种方案：
- ✅ 完全可用
- ✅ 稳定可靠
- ✅ 低延迟
- ✅ 高音质
- ✅ 支持多应用同时捕获（使用多个虚拟设备）

### 代码结构

```
internal/audio/
├── appcapture.go           # 跨平台接口定义
├── appcapture_stub.go      # macOS/Linux 存根实现
└── wasapi_windows.go       # Windows WASAPI 实现 ✅ 新建
    ├── wasapiCaptureImpl   # WASAPI 捕获实现类
    ├── ListApplications()  # 枚举应用 ✅ 已完成
    ├── StartCapture()      # 开始捕获框架
    ├── StopCapture()       # 停止捕获
    ├── captureLoop()       # 捕获循环框架
    └── getProcessName()    # 获取进程名 ✅ 已完成
```

### 编译标签

使用 Go build tags 实现平台分离：

- `wasapi_windows.go`: `// +build windows` - 仅 Windows 编译
- `appcapture_stub.go`: `// +build !windows` - 非 Windows 平台

### GUI 集成

在 `internal/gui/app.go` 中：

```go
// 初始化应用捕获管理器
appCaptureManager, err := audio.NewApplicationCaptureManager()

// 创建应用选择器
appSelect := widget.NewSelect(appNames, func(selected string) {
    // 处理应用选择
})

// 刷新按钮
refreshButton := widget.NewButton("🔄 刷新", func() {
    apps, _ := a.appCaptureManager.ListApplications()
    // 更新下拉列表
})
```

### 测试建议

#### 测试应用列表

推荐测试以下应用：

1. **音乐播放器**
   - Spotify
   - iTunes
   - foobar2000
   - MusicBee

2. **视频播放器**
   - VLC
   - Windows Media Player
   - PotPlayer

3. **浏览器**
   - Google Chrome（播放 YouTube）
   - Microsoft Edge
   - Firefox

4. **通讯软件**
   - Discord
   - Microsoft Teams
   - Zoom
   - Skype

5. **游戏平台**
   - Steam
   - Epic Games

#### 测试步骤

1. 启动测试应用并播放音频
2. 打开 Audio Mixer GUI
3. 点击"刷新"按钮
4. 检查下拉列表中是否显示该应用
5. 验证应用名称、图标、状态是否正确

### 性能注意事项

1. **COM 资源管理**
   - 所有 COM 对象都使用 `defer Release()` 正确释放
   - 避免内存泄漏

2. **线程安全**
   - 使用 `sync.Mutex` 保护共享状态
   - 捕获循环运行在独立 goroutine

3. **错误处理**
   - 枚举失败不会导致程序崩溃
   - 单个应用获取失败不影响其他应用
   - 提供友好的错误提示

### 未来增强方向

如果需要完整的单应用 loopback 捕获，可以考虑：

#### 方案 1: Windows Audio Graph API

使用 Windows 10+ 的 Audio Graph API：
- 需要 C++/WinRT
- 需要 UWP 应用权限
- 可以创建应用到应用的音频路由

**参考资源**:
- [Audio Graphs](https://docs.microsoft.com/en-us/windows/uwp/audio-video-camera/audio-graphs)

#### 方案 2: Audio Endpoint Isolation

使用虚拟音频端点 + WASAPI Session Volume Control：
- 通过控制会话音量实现隔离
- 将非目标应用静音
- 捕获混合输出

**限制**:
- 会影响用户听到的音频
- 不够优雅

#### 方案 3: 继续使用 VB-Cable (推荐)

**优势**:
- 已经完全实现
- 用户配置简单
- 性能优秀
- 稳定可靠

### 相关文档

- 📖 [WASAPI_FEATURE_STATUS.md](WASAPI_FEATURE_STATUS.md) - 功能状态
- 🚀 [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md) - Windows 配置指南
- 🔧 [WINDOWS_QUICK_REFERENCE.md](WINDOWS_QUICK_REFERENCE.md) - 快速参考
- 💡 [SOLUTION_SUMMARY.md](SOLUTION_SUMMARY.md) - 解决方案总结

### 依赖项

```go
// go.mod
require (
    fyne.io/fyne/v2 v2.4.5
    github.com/go-ole/go-ole v1.3.0  // ← 新增
    github.com/gordonklaus/portaudio v0.0.0-20230709114228-aafa478834f5
)
```

### 构建说明

#### Windows

```bash
# 安装依赖
go get -u github.com/go-ole/go-ole

# 编译
go build -ldflags="-H windowsgui" -o audio-mixer-gui.exe cmd/gui/main.go

# 运行
./audio-mixer-gui.exe
```

#### 跨平台编译

```bash
# Windows (from macOS/Linux)
GOOS=windows GOARCH=amd64 go build -o audio-mixer-gui.exe cmd/gui/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o audio-mixer-gui cmd/gui/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o audio-mixer-gui cmd/gui/main.go
```

### 常见问题

#### Q: 为什么应用列表为空？
A: 确保有应用正在播放音频。WASAPI 只能检测到活动的音频会话。

#### Q: 选择应用后没有音频？
A: 当前需要配合 VB-Cable 使用。请参考 WINDOWS_SETUP_GUIDE.md 配置虚拟设备。

#### Q: 编译时找不到 go-ole？
A: 运行 `go get github.com/go-ole/go-ole` 安装依赖。

#### Q: 显示 "COM 初始化失败"？
A: 确保在 Windows 10+ 系统上运行，且有音频设备。

#### Q: macOS 上能用吗？
A: macOS 需要使用 BlackHole 虚拟设备方案。应用枚举需要 ScreenCaptureKit API（待实现）。

### 总结

**当前状态**:
- ✅ Windows 应用枚举完全可用
- ✅ GUI 集成完成
- ✅ 跨平台编译支持
- ⚠️ 音频捕获需要配合虚拟设备

**推荐使用方式**:
1. **Windows**: VB-Cable + WASAPI 应用选择器
2. **macOS**: BlackHole + Multi-Output Device
3. **Linux**: PulseAudio Loopback

**开发时间**:
- WASAPI 枚举实现: ✅ 已完成
- 完整 loopback 捕获: 如需要，预计 2-3 天

---

**最后更新**: 2025-11-27
**状态**: WASAPI 应用枚举功能已完成，推荐配合 VB-Cable 使用
