# Windows WASAPI 实现方案

## 你的需求

> "创建一个新的虚拟设备来输出，而不是选择现有的虚拟设备，并且系统音频是抓取现有的软件源的声音，而不是总的输出，类似于obs的单独捕获音乐软件的声音"

## 技术分析

### 需求拆解

1. **创建虚拟音频设备** - Audio Mixer 自己创建，不依赖 VB-Cable
2. **捕获特定应用音频** - 类似 OBS，从应用列表中选择

### Windows 音频架构

```
┌─────────────────────────────────────┐
│ 应用层 (Spotify, Chrome, etc.)      │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ WASAPI (Windows Audio Session API)  │ ← 我们可以在这里捕获
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ 音频引擎 (Audio Engine)              │
└──────────────┬──────────────────────┘
               │
               ▼
┌─────────────────────────────────────┐
│ 驱动层 (Audio Drivers)               │ ← 创建虚拟设备需要在这里
└─────────────────────────────────────┘
```

## 两个功能的实现难度

### ✅ 功能 1: 捕获特定应用音频 (可实现)

**技术方案**: Windows Audio Session API (WASAPI)

**难度**: ⭐⭐⭐☆☆ (中等)

**实现方法**:
1. 使用 WASAPI 枚举所有音频会话
2. 获取每个会话的进程信息
3. 用户在 GUI 中选择要捕获的应用
4. 使用 Loopback Capture 捕获该应用的音频

**优势**:
- ✅ 不需要安装额外驱动
- ✅ Windows 原生 API
- ✅ 类似 OBS 的功能

**劣势**:
- ⚠️ 需要使用 CGo 调用 Windows COM 接口
- ⚠️ 代码复杂度较高
- ⚠️ 仅支持 Windows 10+

### ❌ 功能 2: 创建虚拟音频设备 (几乎不可行)

**技术方案**: 需要编写音频驱动

**难度**: ⭐⭐⭐⭐⭐ (极难)

**为什么很难**:

1. **需要编写内核驱动**
   - Windows Driver Kit (WDK)
   - 需要数字签名
   - 需要通过 WHQL 认证

2. **需要实现 Audio Processing Object (APO)**
   - C++ 开发
   - COM 接口实现
   - 复杂的音频处理逻辑

3. **用户安装体验差**
   - 需要管理员权限
   - 需要禁用驱动签名验证（测试）
   - 可能被杀毒软件拦截

**为什么 VB-Cable 能做到**:
- VB-Cable 是一个完整的驱动程序项目
- 开发了多年，经过充分测试
- 已经通过了 Windows 驱动签名认证

## 推荐的实际解决方案

### 方案 A: WASAPI 应用捕获 + 使用 VB-Cable 输出（推荐）

#### 实现内容：

1. **为 Audio Mixer 添加应用选择器**
   ```
   GUI 界面：
   ┌────────────────────────────────┐
   │ Input 2 (应用音频):             │
   │ [🎵 Spotify ▼]                  │
   │   - Spotify                     │
   │   - Chrome                      │
   │   - Discord                     │
   │   - VLC Player                  │
   └────────────────────────────────┘
   ```

2. **使用 WASAPI 捕获选定应用的音频**
   - 自动检测正在播放的应用
   - 显示应用图标和名称
   - 实时捕获音频流

3. **输出仍使用 VB-Cable**
   - 不需要重新发明轮子
   - VB-Cable 稳定可靠
   - 一次安装，永久使用

#### 优势：
- ✅ **最佳用户体验** - 类似 OBS，直接选择应用
- ✅ **不需要手动配置** - 自动检测应用
- ✅ **稳定可靠** - 输出使用成熟的 VB-Cable
- ✅ **开发可行** - 可以在 2-3 天内实现

#### 效果：

```
用户操作：
1. 打开 Audio Mixer GUI
2. Input 2 → 从下拉列表选择 "Spotify"
3. Input 1 → 选择麦克风
4. Output → 输入 "CABLE-B Input"
5. 点击 Start Mixer

Audio Mixer 自动：
- 检测 Spotify 进程
- 捕获 Spotify 的音频流
- 与麦克风混音
- 输出到 VB-Cable
```

### 方案 B: 纯软件虚拟设备（实验性）

使用 PortAudio 的 callback 机制创建"软虚拟设备"：

```go
// 创建一个内存中的虚拟设备
type SoftwareVirtualDevice struct {
    buffer     *RingBuffer
    callbacks  []AudioCallback
}

// 其他应用通过 API 连接到这个虚拟设备
func (d *SoftwareVirtualDevice) RegisterConsumer(callback AudioCallback) {
    d.callbacks = append(d.callbacks, callback)
}
```

**问题**:
- 其他应用无法"看到"这个虚拟设备
- 只能通过 Audio Mixer 的 API 使用
- 不适合 OBS/Zoom 等第三方软件

## 我建议的实现计划

### Phase 1: WASAPI 应用捕获（2-3天）

**文件结构**:
```
internal/audio/
├── wasapi_windows.go      # WASAPI 实现（Windows）
├── appcapture.go          # 应用音频捕获接口
└── appcapture_stub.go     # macOS/Linux 存根

internal/gui/
└── app.go                 # 添加应用选择器 UI
```

**功能**:
1. ✅ 列出所有正在播放音频的应用
2. ✅ 显示应用名称和图标
3. ✅ 选择应用后自动捕获其音频
4. ✅ 实时显示捕获状态

### Phase 2: GUI 优化（1天）

**改进**:
1. 应用列表自动刷新
2. 显示应用音量电平
3. 添加应用音频预览
4. 保存常用应用配置

### Phase 3: 输出优化（1天）

**改进**:
1. 自动检测已安装的虚拟设备
2. 如果没有虚拟设备，提示下载 VB-Cable
3. 一键安装 VB-Cable（可选）

## 代码示例（伪代码）

### 应用音频捕获

```go
package audio

// ApplicationCapture 应用音频捕获
type ApplicationCapture struct {
    processID   uint32
    processName string
    callback    func([]float32)
}

// ListAudioApplications 列出所有正在播放音频的应用
func ListAudioApplications() ([]*ApplicationInfo, error) {
    // Windows: 使用 WASAPI
    // macOS: 使用 ScreenCaptureKit (需要 macOS 13+)
    // Linux: 使用 PulseAudio
}

// CaptureApplication 捕获指定应用的音频
func CaptureApplication(app *ApplicationInfo) (*ApplicationCapture, error) {
    // 使用 WASAPI Loopback Capture
    // 只捕获该应用的音频流
}
```

### GUI 集成

```go
// 在 buildDeviceSection 中添加应用选择器
func (a *App) buildInput2Section() fyne.CanvasObject {
    // 获取正在播放音频的应用
    apps, _ := audio.ListAudioApplications()

    appNames := make([]string, len(apps))
    for i, app := range apps {
        appNames[i] = fmt.Sprintf("🎵 %s", app.Name)
    }

    a.appSelect = widget.NewSelect(appNames, func(selected string) {
        // 用户选择了应用
        appID := extractAppID(selected)
        a.captureApp(appID)
    })

    return container.NewVBox(
        widget.NewLabel("Input 2 (应用音频):"),
        a.appSelect,
        widget.NewButton("刷新应用列表", a.refreshApps),
    )
}
```

## 开发工作量估算

| 任务 | 时间 | 难度 |
|------|------|------|
| WASAPI 基础集成 | 1天 | ⭐⭐⭐⭐ |
| 应用枚举功能 | 0.5天 | ⭐⭐⭐ |
| 应用音频捕获 | 1天 | ⭐⭐⭐⭐ |
| GUI 应用选择器 | 0.5天 | ⭐⭐ |
| 测试和调试 | 1天 | ⭐⭐⭐ |
| **总计** | **4天** | **中等** |

## 用户最终体验

### 使用流程：

```
步骤 1: 安装 VB-Cable（一次性，10分钟）
  - 下载 VB-Cable
  - 安装驱动
  - 重启电脑

步骤 2: 启动 Audio Mixer
  - 打开 Audio Mixer GUI

步骤 3: 选择应用（30秒）
  - Input 1: 选择麦克风
  - Input 2: 从列表选择 "🎵 Spotify"  ← 新功能！
  - Output: CABLE-B Input

步骤 4: 开始混音
  - 点击 Start Mixer
  - Audio Mixer 自动捕获 Spotify 音频

步骤 5: 在 OBS 使用
  - OBS 音频源: CABLE-B Output
```

### 对比现有方案：

| 方面 | 当前方案 | 新方案（WASAPI） |
|------|---------|-----------------|
| 应用配置 | 需要在应用内设置输出 | 直接在 Audio Mixer 选择 ✅ |
| 易用性 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 灵活性 | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| 虚拟设备依赖 | 需要 VB-Cable | 仍需要 VB-Cable ⚠️ |

## 关于"创建虚拟设备"的替代方案

由于创建真正的虚拟设备极其困难，我建议：

### 选项 1: 使用 VB-Cable（最实用）
- ✅ 稳定可靠
- ✅ 免费开源
- ✅ 一次安装永久使用
- ✅ 被广泛使用和信任

### 选项 2: 内置 VB-Cable 安装器
我可以为 Audio Mixer 添加一个功能：
- 检测系统是否安装了 VB-Cable
- 如果没有，提供一键下载和安装
- 自动完成配置

```go
func CheckVirtualDevice() error {
    if !isVBCableInstalled() {
        // 显示安装向导
        showInstallWizard()
    }
}
```

### 选项 3: 提供预配置的 VB-Cable 包
- Audio Mixer 安装包包含 VB-Cable
- 安装 Audio Mixer 时自动安装 VB-Cable
- 用户无感知，开箱即用

## 总结与建议

### 我的建议：

**实现 WASAPI 应用捕获功能**，让用户可以：
1. ✅ **直接选择应用** - 不需要在每个应用中配置输出
2. ✅ **类似 OBS 体验** - 从列表中点选即可
3. ✅ **自动捕获** - Audio Mixer 自动处理

**继续使用 VB-Cable 作为输出**，因为：
1. ✅ **创建虚拟设备需要写驱动** - 太复杂，投入产出比低
2. ✅ **VB-Cable 已经很好用** - 稳定、免费、可靠
3. ✅ **可以优化安装体验** - 提供一键安装或内置安装器

### 如果你同意这个方案：

我可以立即开始开发 WASAPI 应用捕获功能，预计 3-4 天完成。

### 你会得到：

```
┌─────────────────────────────────────┐
│ Audio Mixer - 应用音频混音器         │
├─────────────────────────────────────┤
│ Input 1 (麦克风):                    │
│ [MacBook Pro Microphone      ▼]    │
│                                      │
│ Input 2 (应用音频):                  │
│ [🎵 Spotify                  ▼]    │  ← 新功能
│   🎵 Spotify                         │
│   🌐 Google Chrome                   │
│   💬 Discord                         │
│   ▶️ VLC Media Player                │
│                                      │
│ Output:                              │
│ [CABLE-B Input] [检测设备]          │
│                                      │
│ [Start Mixer] [Stop Mixer]          │
└─────────────────────────────────────┘
```

### 下一步：

你想要我：
1. **立即开始实现 WASAPI 应用捕获功能**？
2. **还是继续使用当前的 VB-Cable 手动配置方案**？

请告诉我你的选择！
