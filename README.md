# Audio Mixer

跨平台实时音频混音工具，支持麦克风和系统音频的实时混音，输出到虚拟设备供其他应用使用。

## 功能特性

- ✅ **双路音频输入**:
  - Input 1: 麦克风/线路输入
  - Input 2: 系统音频（需要虚拟音频设备）
- ✅ **虚拟设备输出**: 混音后的音频输出到虚拟设备，可被其他应用使用
- ✅ **应用音频捕获** (Windows):
  - 🆕 列出所有正在播放音频的应用程序
  - 🆕 显示应用名称、进程信息、播放状态
  - 🆕 类似 OBS 的应用选择功能
  - 📖 详见 [WASAPI_IMPLEMENTATION_NOTES.md](WASAPI_IMPLEMENTATION_NOTES.md)
- ✅ **实时混音**: 低延迟(<30ms)音频混音处理
- ✅ **独立音量控制**: 每路输入独立调节(0-200%)
- ✅ **软削波防护**: 防止音频爆音和失真
- ✅ **实时监控**: 显示各路音频电平和处理延迟
- ✅ **配置持久化**: 自动保存设备和音量设置
- ✅ **跨平台支持**: macOS优先,支持Windows和Linux
- ✅ **自动检测**: 自动检测并配置虚拟音频设备

## 技术栈

- **语言**: Go 1.21+
- **音频库**: PortAudio (github.com/gordonklaus/portaudio)
- **GUI框架**: Fyne v2 (跨平台图形界面)
- **Windows API**: go-ole (COM 接口，用于 WASAPI)
- **架构**: 低延迟实时音频处理,零GC暂停设计

## 系统要求

### macOS
- macOS 10.12+
- 已安装 PortAudio: `brew install portaudio`
- 麦克风权限(首次运行时系统会提示)
- **必需**: BlackHole 虚拟音频设备
  - 安装: `brew install blackhole-2ch`
  - 用于捕获系统音频和虚拟输出
  - [详细设置指南](VIRTUAL_DEVICE_SETUP.md)

### Windows
- Windows 10+
- PortAudio库(构建时自动链接)
- **必需**: VB-Cable 或 Virtual Audio Cable
  - 下载: [VB-Audio Virtual Cable](https://vb-audio.com/Cable/)
  - 用于系统音频捕获和虚拟输出
  - [详细设置指南](VIRTUAL_DEVICE_SETUP.md)
- **新功能**: WASAPI 应用音频枚举
  - 自动列出正在播放音频的应用
  - 需要 go-ole 库（构建时自动安装）
  - 📖 [WASAPI 实现说明](WASAPI_IMPLEMENTATION_NOTES.md)

### Linux
- PulseAudio或ALSA
- PortAudio库: `sudo apt-get install portaudio19-dev`

## 安装

### 快速构建

**推荐方式** - 使用构建脚本：

```bash
# macOS/Linux
./build.sh

# Windows (PowerShell)
.\build.ps1

# 或使用 Makefile
make gui
```

📖 **详细构建指南**: [BUILD_QUICK_REFERENCE.md](BUILD_QUICK_REFERENCE.md)

---

### 从源码构建

1. 安装依赖:

**macOS:**
```bash
brew install portaudio
```

**Linux (Debian/Ubuntu):**
```bash
sudo apt-get install portaudio19-dev
```

**Windows:**
- PortAudio 会自动通过 CGO 链接
- 可能需要安装 MinGW-w64 (C 编译器)
- WASAPI 支持需要 `github.com/go-ole/go-ole` (自动下载)

2. 克隆仓库:
```bash
git clone https://github.com/entropy/audio-mixer.git
cd audio-mixer
```

3. 构建:

**CLI 版本:**
```bash
go build -o audio-mixer .
```

**GUI 版本:**
```bash
# macOS/Linux
go build -o audio-mixer-gui ./cmd/gui

# Windows (隐藏控制台窗口)
go build -ldflags="-H windowsgui" -o audio-mixer-gui.exe ./cmd/gui
```

4. 运行:
```bash
./audio-mixer        # CLI
./audio-mixer-gui    # GUI
```

📖 **Windows 详细编译**: [BUILD_WINDOWS.md](BUILD_WINDOWS.md)

## 使用指南

### 第一步：安装虚拟音频设备

**这是必需的步骤！** 否则无法捕获系统音频。

#### macOS:
```bash
# 安装 BlackHole
brew install blackhole-2ch

# 创建 Multi-Output Device (在 Audio MIDI Setup 中)
# 详见: VIRTUAL_DEVICE_SETUP.md
```

#### Windows:
```bash
# 下载并安装 VB-Cable
# https://vb-audio.com/Cable/
```

#### 完整设置指南:

**macOS 用户:**
- 🚀 **[快速配置指南 (QUICK_SETUP_GUIDE.md)](QUICK_SETUP_GUIDE.md)** - **5分钟快速配置，捕获特定应用音频**

**Windows 用户:**
- 🪟 **[Windows 配置指南 (WINDOWS_SETUP_GUIDE.md)](WINDOWS_SETUP_GUIDE.md)** - **10分钟 Windows 完整配置**

**通用文档:**
- 📖 **[虚拟设备设置指南 (VIRTUAL_DEVICE_SETUP.md)](VIRTUAL_DEVICE_SETUP.md)** - 跨平台虚拟设备安装
- 🔧 **[高级音频路由 (ADVANCED_AUDIO_ROUTING.md)](ADVANCED_AUDIO_ROUTING.md)** - 专业音频路由方案
- 💡 **[解决方案总结 (SOLUTION_SUMMARY.md)](SOLUTION_SUMMARY.md)** - 技术原理和方案对比

---

### 第二步：运行 Audio Mixer

#### GUI模式 (推荐)

1. 构建GUI版本:
```bash
go build -o audio-mixer-gui ./cmd/gui
```

2. 运行GUI:
```bash
./audio-mixer-gui
```

3. 配置设备:
   - **Input 1 (麦克风)**: 选择你的麦克风设备
   - **Input 2 (系统音频)**: 选择虚拟设备或 "<Auto Detect Loopback>"
     - 自动检测会查找 BlackHole、Soundflower 等虚拟设备
   - **Output (虚拟输出设备名称)**: 输入虚拟设备名称
     - 例如: `BlackHole 2ch`, `BlackHole 16ch`, `VB-Cable`
     - 点击 "检测设备" 按钮验证设备是否存在
     - 📖 **[自定义输出设备指南](CUSTOM_OUTPUT_DEVICE.md)**
   - **应用音频捕获 (Windows)**:
     - 点击 "🔄 刷新" 按钮查看正在播放音频的应用
     - 选择要捕获的应用（需配合 VB-Cable 使用）
     - 📖 **[WASAPI 功能说明](WASAPI_IMPLEMENTATION_NOTES.md)**

4. 调节音量:
   - 使用滑块调节 Input 1（麦克风）音量(0-200%)
   - 使用滑块调节 Input 2（系统音频）音量(0-200%)
   - 使用滑块调节 Master（总输出）音量(0-200%)

5. 开始混音:
   - 点击 "Start Mixer" 开始混音
   - 实时查看音频电平表
   - 点击 "Stop Mixer" 停止

6. 自定义字体(可选,用于更好的中文显示):
```bash
# 使用自定义字体文件 (仅支持 TTF/OTF,不支持 TTC)
./audio-mixer-gui -font /path/to/your/font.ttf

# macOS 使用 Arial Unicode (系统自带,支持中文)
./audio-mixer-gui -font /System/Library/Fonts/Supplemental/Arial\ Unicode.ttf

# 或设置环境变量
FYNE_FONT=/System/Library/Fonts/Supplemental/Arial\ Unicode.ttf ./audio-mixer-gui
```

### CLI模式

1. 启动程序:
```bash
./audio-mixer
```

2. 程序会列出所有可用的音频设备(输入和输出)

3. 按提示选择设备:
   - **Input 1**: 麦克风设备(输入 `-1` 使用默认设备)
   - **Input 2**: 应用音频设备(输入 `-2` 跳过第二输入)
   - **Output**: 输出设备(可以是虚拟设备)

4. 设置音量:
   - 范围: 0.0 - 2.0 (0% - 200%)
   - 1.0 = 100% (原始音量)
   - 大于1.0 可以增益信号

5. 混音器开始运行,实时显示:
   - 各路输入电平(dB)
   - 输出电平
   - 处理延迟

6. 按 `Ctrl+C` 停止

### 配置文件

配置自动保存在: `~/.audio-mixer/config.json`

可以手动编辑此文件修改设置:

```json
{
  "sample_rate": 48000,
  "buffer_size": 512,
  "channels": 2,
  "input1_device_index": -1,
  "input2_device_index": -1,
  "output_device_index": -1,
  "input1_gain": 1.0,
  "input2_gain": 1.0,
  "master_gain": 1.0
}
```

### 虚拟音频设备

为了捕获应用程序音频和输出到语音软件,需要虚拟音频设备:

#### macOS - BlackHole

1. 下载安装 [BlackHole](https://github.com/ExistentialAudio/BlackHole)
```bash
brew install blackhole-2ch
```

2. 在"音频MIDI设置"中创建多输出设备:
   - 打开 `/Applications/Utilities/Audio MIDI Setup.app`
   - 点击左下角 `+` → "创建多输出设备"
   - 勾选你的实际输出设备和BlackHole
   - 将应用程序音频输出到BlackHole

3. 在Audio Mixer中:
   - Input 2 选择 BlackHole(捕获应用音频)
   - Output 选择另一个BlackHole或虚拟设备

#### Windows - VB-Cable

1. 下载 [VB-Cable](https://vb-audio.com/Cable/)
2. 安装虚拟音频驱动
3. 设置应用程序输出到VB-Cable Input
4. Audio Mixer的Input 2选择VB-Cable Output

#### Linux - PulseAudio

使用PulseAudio的monitor设备:
```bash
pacmd list-sources | grep -e 'name:' -e 'device.description'
```

选择想监控的应用的monitor设备作为Input 2

## 使用场景

### 场景1: 在Discord中播放音乐同时说话

1. 安装BlackHole虚拟设备
2. 音乐播放器输出到BlackHole
3. Audio Mixer:
   - Input 1: 麦克风
   - Input 2: BlackHole(音乐)
   - Output: 另一个虚拟设备
4. Discord输入选择Audio Mixer的输出设备

### 场景2: 游戏直播混音

1. Input 1: 麦克风(你的声音)
2. Input 2: 游戏音频(通过虚拟设备)
3. Output: OBS虚拟输入
4. 独立调节麦克风和游戏音量

### 场景3: 播客录制

1. Input 1: 主持人麦克风
2. Input 2: 音乐/音效(从播放器)
3. Output: 录音软件输入
4. 实时监控音量电平

## 性能指标

- **延迟**: 通常 < 20ms
- **CPU占用**: < 5% (单核)
- **内存**: < 50MB
- **采样率**: 48000 Hz
- **位深度**: 32-bit float
- **缓冲区**: 512帧(可配置)

## 项目结构

```
audio-mixer/
├── cmd/
│   └── gui/
│       └── main.go            # GUI 程序入口
├── internal/
│   ├── audio/
│   │   ├── mixer.go           # 核心混音引擎
│   │   ├── device.go          # 设备管理
│   │   ├── loopback.go        # 虚拟设备检测
│   │   ├── appcapture.go      # 应用捕获接口 (跨平台)
│   │   ├── wasapi_windows.go  # Windows WASAPI 实现 🆕
│   │   └── appcapture_stub.go # macOS/Linux 存根
│   ├── config/
│   │   └── config.go          # 配置管理
│   └── gui/
│       ├── app.go             # GUI 主界面
│       └── font.go            # 字体管理
├── go.mod
├── go.sum
├── README.md
├── WASAPI_IMPLEMENTATION_NOTES.md  # WASAPI 实现说明 🆕
├── WASAPI_FEATURE_STATUS.md        # WASAPI 功能状态 🆕
├── WINDOWS_SETUP_GUIDE.md          # Windows 配置指南
└── QUICK_SETUP_GUIDE.md            # 快速配置指南
```

## 架构设计

### 音频处理流程

```
Input 1 (Mic) ──→ [Callback] ──→ [Ring Buffer 1] ──┐
                                                      │
Input 2 (App) ──→ [Callback] ──→ [Ring Buffer 2] ──┤
                                                      ├──→ [Mixer] ──→ [Output]
                                                      │     ↓
Gain Controls ────────────────────────────────────────→ [Volume]
                                                            ↓
                                                      [Soft Clip]
```

### 混音算法

```go
// 简单加法混音 + 软削波
output[i] = (input1[i] * gain1 + input2[i] * gain2) * masterGain

// 软削波(防止失真)
if output[i] > 0.9 {
    output[i] = tanh_soft_clip(output[i])
}
```

### 性能优化

1. **sync.Pool**: 减少内存分配
2. **Ring Buffer**: 无锁缓冲区设计
3. **Atomic操作**: 线程安全的音量控制
4. **零拷贝**: 直接操作音频buffer
5. **低GC**: 避免运行时暂停

## 故障排除

### 问题: 没有声音

- 检查设备选择是否正确
- 确认音量不是0
- 检查虚拟设备是否正确安装
- 查看实时电平是否有信号

### 问题: 延迟过高

- 减小buffer size(在config.json中设置)
- 关闭不必要的音频处理
- 检查CPU占用

### 问题: 音质失真

- 降低增益(Gain < 1.0)
- 检查输出电平是否过高(>0dB)
- 确保输入信号不削波

### 问题: 无法找到设备(macOS)

- 授予麦克风权限: 系统偏好设置 → 安全性与隐私 → 麦克风
- 重启程序
- 检查虚拟设备是否正确安装

### 问题: GUI中文显示为方块

这是字体不支持中文字符导致的,可以通过以下方式解决:

**重要提示: Fyne GUI框架仅支持 TTF/OTF 单字体文件,不支持 TTC (字体集合) 文件**

**方法1: 使用 GUI 内置字体选择器**
- 启动程序后,在 "Font Settings" 选项卡中选择可用字体
- 程序会自动检测系统中可用的中文字体

**方法2: 使用命令行参数指定字体**
```bash
# macOS - 使用系统自带的 Arial Unicode (支持中文)
./audio-mixer-gui -font /System/Library/Fonts/Supplemental/Arial\ Unicode.ttf

# 使用用户安装的字体
./audio-mixer-gui -font ~/Library/Fonts/SourceHanSansSC-Regular.otf
```

**方法3: 设置环境变量**
```bash
export FYNE_FONT=/System/Library/Fonts/Supplemental/Arial\ Unicode.ttf
./audio-mixer-gui
```

**方法4: 下载并使用开源中文字体**
```bash
# 下载思源黑体 (Source Han Sans)
wget https://github.com/adobe-fonts/source-han-sans/releases/download/2.004R/SourceHanSansSC.zip
unzip SourceHanSansSC.zip
./audio-mixer-gui -font ./SourceHanSansSC/OTF/SimplifiedChinese/SourceHanSansSC-Regular.otf

# 或下载 Noto Sans CJK
wget https://github.com/googlefonts/noto-cjk/releases/download/Sans2.004/NotoSansCJKsc.zip
unzip NotoSansCJKsc.zip
./audio-mixer-gui -font ./NotoSansCJKsc-Regular.otf
```

**可用的系统字体路径(macOS):**
- Arial Unicode MS (推荐): `/System/Library/Fonts/Supplemental/Arial Unicode.ttf`
- ⚠️ 注意: macOS 系统字体如 PingFang.ttc、STHeiti.ttc 是 TTC 格式,Fyne 不支持
- 需要从 TTC 文件中提取单个字体,或下载上述开源字体

## 开发计划

### Phase 1 ✅ (已完成)
- [x] 基础音频设备枚举
- [x] 单路音频输入/输出
- [x] 双路混音功能
- [x] 音量控制
- [x] 配置管理
- [x] CLI界面

### Phase 2 ✅ (已完成)
- [x] GUI界面(Fyne)
- [x] 可视化音频电平(实时进度条)
- [ ] 音频效果器(EQ、压缩器)
- [ ] 录音功能
- [ ] 预设管理

### Phase 3 (未来)
- [ ] 多路输入混音(3+路)
- [ ] 音频波形/频谱显示
- [ ] 热键支持
- [ ] 插件系统

## 贡献

欢迎提交Issue和Pull Request!

## 许可证

MIT License

## 相关资源

- [PortAudio文档](http://www.portaudio.com/docs/v19-doxydocs/)
- [BlackHole虚拟音频](https://github.com/ExistentialAudio/BlackHole)
- [VB-Cable](https://vb-audio.com/Cable/)
- [Go PortAudio绑定](https://github.com/gordonklaus/portaudio)

## 致谢

- PortAudio团队
- Go社区
- 所有贡献者
