# 高级音频路由方案

## 需求分析

您希望实现：
1. **创建专用虚拟输出设备** - 用于 Audio Mixer 的混音输出
2. **捕获特定应用音频** - 类似 OBS，捕获特定软件（如音乐播放器）的声音

## 技术限制

### macOS
- **无法动态创建虚拟设备**：需要安装内核驱动或 Audio Server Plugin
- **无法直接捕获应用音频**：需要 ScreenCaptureKit API (macOS 13+) 或第三方工具

### 可行方案

## 方案 A: 使用多个 BlackHole 实例（推荐）

### 1. 安装多个 BlackHole 设备

```bash
# 安装 2 通道版本（用于特定应用音频）
brew install blackhole-2ch

# 安装 16 通道版本（用于混音输出）
brew install blackhole-16ch
```

安装后你会有：
- **BlackHole 2ch** - 用于捕获特定应用音频
- **BlackHole 16ch** - 用于 Audio Mixer 输出

### 2. 配置特定应用音频路由

#### 方法 1: 使用 Multi-Output Device（免费）

1. 打开 **Audio MIDI Setup**

2. 创建 **Aggregate Device**（聚合设备）：
   - 点击 `+` → **Create Aggregate Device**
   - 勾选：
     - ✅ BlackHole 2ch
     - ✅ 你的扬声器（可选，用于监听）

3. 在音乐播放器或目标应用中：
   - 设置音频输出为 **BlackHole 2ch**
   - 或设置为刚创建的 Aggregate Device

4. 在 Audio Mixer 中配置：
   ```
   Input 1 (麦克风): MacBook Pro Microphone
   Input 2 (应用音频): BlackHole 2ch  ← 捕获音乐播放器
   Output: BlackHole 16ch  ← 混音后输出
   ```

5. 在 OBS 或其他软件中：
   - 音频源选择 **BlackHole 16ch**

#### 方法 2: 使用 Loopback（商业软件，$99）

Loopback 提供图形化界面，可以：
- 捕获特定应用的音频
- 创建虚拟音频设备
- 灵活路由音频流

下载：https://rogueamoeba.com/loopback/

### 3. 音频流图示

```
┌─────────────────────┐
│ 音乐播放器 (Spotify)│
│ 输出: BlackHole 2ch │
└──────────┬──────────┘
           │
           ▼
┌──────────────────────────────┐
│ Audio Mixer                  │
│ ┌──────────────────────────┐ │
│ │ Input 1: 麦克风          │ │
│ │ Input 2: BlackHole 2ch   │ │◄─ 捕获音乐
│ │         (音乐)            │ │
│ │                           │ │
│ │ 混音处理                 │ │
│ │                           │ │
│ │ Output: BlackHole 16ch   │ │
│ └──────────────────────────┘ │
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│ OBS / 录音软件                │
│ 音频源: BlackHole 16ch        │
└──────────────────────────────┘
```

## 方案 B: 使用 ScreenCaptureKit API（需要开发）

这需要为 Audio Mixer 添加应用音频捕获功能。

### 实现步骤：

1. **添加 ScreenCaptureKit 依赖**（macOS 13+）
2. **创建应用选择器 UI**
3. **捕获选定应用的音频流**
4. **将捕获的音频作为 Input 2**

### 优势：
- ✅ 无需配置虚拟设备
- ✅ 直接捕获应用音频
- ✅ 类似 OBS 的体验

### 劣势：
- ❌ 仅支持 macOS 13+
- ❌ 需要屏幕录制权限
- ❌ 开发工作量大
- ❌ 需要 CGo 调用 Objective-C API

### 代码示例（伪代码）：

```go
// 需要使用 CGo 调用 macOS API
/*
#cgo LDFLAGS: -framework ScreenCaptureKit -framework AVFoundation
#import <ScreenCaptureKit/ScreenCaptureKit.h>
#import <AVFoundation/AVFoundation.h>

// 捕获应用音频
void captureApplicationAudio(int processID, void* callback);
*/
import "C"

type ApplicationAudioCapture struct {
    processID int
    stream    chan []float32
}

func (a *ApplicationAudioCapture) Start() error {
    // 使用 ScreenCaptureKit 捕获应用音频
    // C.captureApplicationAudio(C.int(a.processID), callback)
}
```

## 方案 C: 集成 JACK Audio（跨平台）

JACK Audio Connection Kit 提供专业音频路由功能。

### 安装：

```bash
# macOS
brew install jack

# Linux
sudo apt install jackd2
```

### 配置：

1. 启动 JACK 服务
2. 使用 QjackCtl 图形化连接音频流
3. Audio Mixer 作为 JACK 客户端

### 优势：
- ✅ 专业音频路由
- ✅ 低延迟
- ✅ 跨平台

### 劣势：
- ❌ 学习曲线陡峭
- ❌ 需要额外安装
- ❌ macOS 上不够稳定

## 推荐实施方案

### 短期方案（立即可用）：

**使用多个 BlackHole 实例 + Multi-Output Device**

1. 安装 BlackHole 2ch 和 16ch
2. 配置音乐播放器输出到 BlackHole 2ch
3. Audio Mixer 配置：
   - Input 2: BlackHole 2ch
   - Output: BlackHole 16ch
4. OBS 从 BlackHole 16ch 捕获

### 中期方案（需要开发）：

**添加应用音频捕获功能**

创建新的功能模块：
- `internal/audio/appcapture.go` - 应用音频捕获
- 使用 ScreenCaptureKit API（macOS 13+）
- 在 GUI 中添加应用选择器

### 长期方案（完整解决方案）：

**开发专用音频路由解决方案**

类似 Loopback 的功能：
- 虚拟设备管理
- 应用音频捕获
- 图形化音频路由
- 跨平台支持

## 配置示例

### 示例 1: 捕获 Spotify + 麦克风

```
1. Spotify 设置:
   输出设备: BlackHole 2ch

2. Audio Mixer 配置:
   Input 1: MacBook Pro Microphone
   Input 2: BlackHole 2ch (自动检测)
   Output: BlackHole 16ch

3. OBS 配置:
   音频源: BlackHole 16ch
```

### 示例 2: 捕获游戏音频 + Discord

```
1. 游戏设置:
   音频输出: BlackHole 2ch

2. Discord 设置:
   输出设备: 默认（扬声器）

3. 创建 Multi-Output Device:
   - BlackHole 2ch (游戏音频)
   - BlackHole 8ch (Discord 音频)

4. Audio Mixer 配置:
   Input 1: 麦克风
   Input 2: Multi-Output Device
   Output: BlackHole 16ch
```

## 故障排除

### 问题: 听不到音乐播放器的声音

**原因**: 音频只输出到 BlackHole，没有输出到扬声器

**解决方案**:
1. 创建 Multi-Output Device
2. 同时勾选 BlackHole 和扬声器
3. 设置音乐播放器输出到 Multi-Output Device

### 问题: 延迟太高

**解决方案**:
1. 降低 Audio Mixer 的 Buffer Size (256 samples)
2. 关闭不必要的音频应用
3. 使用有线连接而非蓝牙

### 问题: 音频断断续续

**解决方案**:
1. 提高 Buffer Size (1024 samples)
2. 检查 CPU 使用率
3. 关闭其他高负载应用

## 未来开发计划

如果需要实现类似 OBS 的应用音频捕获功能，我可以：

1. **添加 ScreenCaptureKit 集成**
   - 创建应用选择器 UI
   - 实现音频捕获回调
   - 集成到现有混音器

2. **创建虚拟设备管理器**
   - 检测已安装的虚拟设备
   - 提供设备切换功能
   - 自动配置音频路由

3. **开发图形化路由器**
   - 可视化音频流
   - 拖拽式连接
   - 预设配置管理

请告诉我您希望采用哪个方案，我可以为您实现相关功能！
