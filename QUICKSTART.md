# 快速开始指南

5分钟快速上手Audio Mixer!

## 前提条件

确保你已安装:
- Go 1.21+ (`go version`)
- PortAudio (`brew install portaudio` on macOS)

## 快速安装

```bash
# 1. 进入项目目录
cd audio-mixer

# 2. 下载依赖
go mod download

# 3. 构建程序
make build
# 或者: go build -o audio-mixer .

# 4. 运行
./audio-mixer
```

## 首次使用

### 步骤1: 查看设备列表

程序启动后会自动显示所有可用的音频设备:

```
Available Audio Devices:
------------------------

Input Devices:
  [0] Built-in Microphone (Channels: 1, SR: 48000 Hz) [DEFAULT]
  [1] BlackHole 2ch (Channels: 2, SR: 48000 Hz)

Output Devices:
  [2] Built-in Output (Channels: 2, SR: 48000 Hz) [DEFAULT]
  [3] BlackHole 2ch (Channels: 2, SR: 48000 Hz)
```

### 步骤2: 选择设备

根据提示输入设备编号:

```
Select Input 1 device (Microphone) [current: -1, -1 for default]: 0
Select Input 2 device (Application Audio) [current: -1, -2 to skip]: -2
Select Output device [current: -1, -1 for default]: 2
```

**提示**:
- 输入 `-1` 使用默认设备
- 输入 `-2` 跳过第二输入(只使用麦克风)
- 输入设备编号选择特定设备

### 步骤3: 设置音量

```
Input 1 Gain [current: 1.00]: 1.0
Input 2 Gain [current: 1.00]: 0.3
Master Gain [current: 1.00]: 1.0
```

**音量范围**: 0.0 - 2.0
- 0.0 = 静音
- 1.0 = 100% (推荐)
- 2.0 = 200% (增益)

### 步骤4: 开始混音

配置完成后自动开始:

```
=== Starting Audio Mixer ===
Mixer started successfully!

Press Ctrl+C to stop

Real-time Monitoring:
---------------------
[Input1:  -20.5 dB ████████████░░░░░░░░] [Input2:  -inf dB ░░░░░░░░░░░░░░░░░░░░] [Output:  -20.5 dB ████████████░░░░░░░░] [Latency: 10.234ms]
```

### 步骤5: 停止

按 `Ctrl+C` 优雅停止:

```
^C
Shutting down...
Goodbye!
```

## 常见使用场景

### 场景1: 简单测试麦克风

```bash
./audio-mixer
# Input 1: 选择麦克风(或-1)
# Input 2: 输入-2跳过
# Output: 选择输出设备(或-1)
# 所有Gain保持1.0
```

说话,你应该能看到Input1的电平条跳动!

### 场景2: 在Discord播放音乐

**准备工作**:
1. 安装BlackHole: `brew install blackhole-2ch`
2. 音乐播放器输出到BlackHole
3. Discord输入设置为Audio Mixer的输出

**配置**:
```bash
./audio-mixer
# Input 1: 麦克风(如: 0)
# Input 2: BlackHole(如: 1)
# Output: 另一个虚拟设备或BlackHole
# Input 1 Gain: 1.0
# Input 2 Gain: 0.3  # 音乐音量30%
# Master Gain: 1.0
```

### 场景3: 游戏直播

```bash
./audio-mixer
# Input 1: 麦克风
# Input 2: 游戏音频(通过虚拟设备)
# Output: OBS音频输入
# Input 1 Gain: 1.2  # 提升语音
# Input 2 Gain: 0.8  # 游戏音频
# Master Gain: 0.9
```

## 配置文件

首次运行后,配置保存在:
```
~/.audio-mixer/config.json
```

下次运行会自动加载上次的设置,直接按回车使用上次的值。

### 手动编辑配置

```bash
# 编辑配置
nano ~/.audio-mixer/config.json

# 示例配置
{
  "sample_rate": 48000,
  "buffer_size": 512,
  "channels": 2,
  "input1_device_index": 0,
  "input2_device_index": -2,
  "output_device_index": 2,
  "input1_gain": 1.0,
  "input2_gain": 0.3,
  "master_gain": 1.0
}
```

## Makefile命令

```bash
make build          # 构建程序
make run           # 构建并运行
make clean         # 清理构建文件
make test          # 运行测试
make deps          # 下载Go依赖
make deps-macos    # 安装系统依赖(macOS)
make build-all     # 构建所有平台版本
```

## 故障排除

### 问题: 命令找不到

```bash
# 确保Go已安装
go version

# 确保PortAudio已安装(macOS)
brew list portaudio
```

### 问题: 权限错误

**macOS**: 授予麦克风权限
- 系统偏好设置 → 安全性与隐私 → 隐私 → 麦克风
- 勾选你的终端应用

### 问题: 没有声音

检查:
1. 设备选择是否正确
2. 音量是否为0
3. 输出设备是否正确连接
4. 查看电平条是否有信号

### 问题: 延迟太高

编辑配置文件,减小buffer_size:
```json
{
  "buffer_size": 256
}
```

## 进阶使用

### 使用示例代码

```bash
# 运行简单直通示例
go run examples/simple_passthrough.go
```

### 查看详细文档

- [README.md](README.md) - 完整使用手册
- [INSTALL.md](INSTALL.md) - 详细安装说明
- [EXAMPLES.md](EXAMPLES.md) - 更多使用场景
- [ARCHITECTURE.md](ARCHITECTURE.md) - 技术架构

## 下一步

1. ✅ 尝试基本的音频直通
2. ✅ 配置虚拟音频设备
3. ✅ 在实际场景中使用(Discord/游戏/录音)
4. ✅ 阅读EXAMPLES.md了解更多用法
5. 🔜 等待GUI版本发布

## 需要帮助?

- 📖 查看文档: [README.md](README.md)
- 🐛 报告问题: GitHub Issues
- 💡 功能建议: GitHub Issues

---

**祝你使用愉快!** 🎵🎙️🎧
