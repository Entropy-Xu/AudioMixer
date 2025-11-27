# Audio Mixer

跨平台实时音频混音工具,支持麦克风和应用程序音频的实时混音输出。

## 功能特性

- ✅ **双路音频输入**: 同时捕获麦克风和应用程序音频
- ✅ **实时混音**: 低延迟(<30ms)音频混音处理
- ✅ **独立音量控制**: 每路输入独立调节(0-200%)
- ✅ **软削波防护**: 防止音频爆音和失真
- ✅ **实时监控**: 显示各路音频电平和处理延迟
- ✅ **配置持久化**: 自动保存设备和音量设置
- ✅ **跨平台支持**: macOS优先,支持Windows和Linux

## 技术栈

- **语言**: Go 1.21+
- **音频库**: PortAudio (github.com/gordonklaus/portaudio)
- **架构**: 低延迟实时音频处理,零GC暂停设计

## 系统要求

### macOS
- macOS 10.12+
- 已安装 PortAudio: `brew install portaudio`
- 麦克风权限(首次运行时系统会提示)
- 可选: BlackHole虚拟音频设备(用于应用音频捕获)

### Windows
- Windows 10+
- PortAudio库(构建时自动链接)
- 可选: VB-Cable虚拟音频驱动

### Linux
- PulseAudio或ALSA
- PortAudio库: `sudo apt-get install portaudio19-dev`

## 安装

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
PortAudio会自动通过CGO链接

2. 克隆仓库:
```bash
git clone https://github.com/entropy/audio-mixer.git
cd audio-mixer
```

3. 构建:
```bash
go build -o audio-mixer .
```

4. 运行:
```bash
./audio-mixer
```

## 使用指南

### 基本使用

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
├── main.go                     # 程序入口和CLI界面
├── internal/
│   ├── audio/
│   │   ├── mixer.go           # 核心混音引擎
│   │   ├── device.go          # 设备管理
│   │   └── buffer.go          # 缓冲区管理
│   └── config/
│       └── config.go          # 配置管理
├── go.mod
└── README.md
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

## 开发计划

### Phase 1 ✅ (已完成)
- [x] 基础音频设备枚举
- [x] 单路音频输入/输出
- [x] 双路混音功能
- [x] 音量控制
- [x] 配置管理
- [x] CLI界面

### Phase 2 (规划中)
- [ ] GUI界面(Wails或Fyne)
- [ ] 可视化音频电平(VU表)
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
