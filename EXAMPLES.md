# 使用示例

本文档提供了Audio Mixer的各种使用场景和配置示例。

## 示例1: 基本音频直通

最简单的用例 - 将麦克风直接传递到输出设备。

### 运行
```bash
go run examples/simple_passthrough.go
```

### 用途
- 测试音频设备是否工作
- 验证延迟
- 测试音频质量

### 代码说明
```go
config := audio.DefaultMixerConfig()
config.Input1Device = micDevice    // 麦克风
config.Input2Device = nil          // 不使用第二输入
config.OutputDevice = outputDevice
config.Input1Gain = 1.0           // 100%
```

## 示例2: 在Discord中播放音乐

在语音通话中同时说话和播放音乐。

### 设置步骤

1. **安装BlackHole** (macOS):
```bash
brew install blackhole-2ch
```

2. **配置音频路由**:
   - 打开"音频MIDI设置"
   - 创建"多输出设备",包含:
     - BlackHole 2ch
     - 你的扬声器(可选,用于监听)
   - 将音乐播放器的输出设置为此多输出设备

3. **运行Audio Mixer**:
```bash
./audio-mixer
```

4. **配置设备**:
   - Input 1: 选择你的麦克风
   - Input 2: 选择"BlackHole 2ch"
   - Output: 选择另一个BlackHole或虚拟设备

5. **Discord设置**:
   - 输入设备: 选择Audio Mixer的输出设备

### 配置示例
```json
{
  "input1_device_index": 0,     // 麦克风
  "input2_device_index": 5,     // BlackHole(捕获音乐)
  "output_device_index": 6,     // 虚拟输出
  "input1_gain": 1.0,           // 语音100%
  "input2_gain": 0.3,           // 音乐30%
  "master_gain": 1.0
}
```

### 音量调节建议
- **语音**: 保持100% (1.0)
- **音乐**: 20-40% (0.2-0.4)
- **Master**: 根据需要调整

## 示例3: 游戏直播混音

将游戏音频和麦克风混合,输出到OBS。

### Windows设置

1. **安装VB-Cable**:
   - 下载并安装VB-Cable
   - 重启电脑

2. **游戏设置**:
   - 音频输出: CABLE Input

3. **OBS设置**:
   - 音频输入: Audio Mixer的输出设备

4. **Audio Mixer配置**:
   ```json
   {
     "input1_device_index": 0,    // 麦克风
     "input2_device_index": 2,    // CABLE Output(游戏音频)
     "output_device_index": 3,    // 另一个虚拟设备
     "input1_gain": 1.2,          // 提升麦克风音量
     "input2_gain": 0.8,          // 降低游戏音量
     "master_gain": 1.0
   }
   ```

### macOS设置

使用BlackHole类似方式:
```json
{
  "input1_gain": 1.5,    // 麦克风增益
  "input2_gain": 0.6,    // 游戏音量
  "master_gain": 0.9
}
```

## 示例4: 播客录制

录制包含背景音乐的播客。

### 设置

1. **设备配置**:
   - Input 1: 主持人麦克风
   - Input 2: 音乐播放器(通过虚拟设备)
   - Output: 录音软件输入(如Audacity)

2. **配置文件**:
```json
{
  "sample_rate": 48000,
  "buffer_size": 512,
  "input1_device_index": 0,
  "input2_device_index": 4,
  "output_device_index": 5,
  "input1_gain": 1.0,
  "input2_gain": 0.15,    // 背景音乐很轻
  "master_gain": 1.0
}
```

3. **使用技巧**:
   - 背景音乐保持在15-25%
   - 监控电平,避免超过-6dB
   - 使用优质麦克风

## 示例5: 在线教学

教师讲课同时播放演示视频音频。

### 配置

```json
{
  "input1_device_index": 0,      // 教师麦克风
  "input2_device_index": 3,      // 演示视频音频
  "output_device_index": 4,      // 会议软件输入
  "input1_gain": 1.0,            // 讲师声音清晰
  "input2_gain": 0.5,            // 演示音频适中
  "master_gain": 1.0
}
```

### 使用场景
- Zoom演示
- Teams会议
- Google Meet讲课

## 示例6: 多设备测试

快速测试不同设备组合。

### 脚本
```bash
#!/bin/bash

# 列出所有设备
./audio-mixer << EOF





EOF

# 测试默认设备
./audio-mixer << EOF
-1
-2
-1



EOF
```

## 示例7: 高音质配置

追求最佳音质的配置。

### 配置
```json
{
  "sample_rate": 96000,      // 高采样率
  "buffer_size": 1024,       // 较大缓冲(更稳定)
  "channels": 2,
  "input1_gain": 1.0,
  "input2_gain": 1.0,
  "master_gain": 0.95        // 留些余量
}
```

### 使用场景
- 音乐制作
- 专业录音
- 高保真回放

## 示例8: 低延迟配置

追求最低延迟的配置。

### 配置
```json
{
  "sample_rate": 48000,
  "buffer_size": 256,        // 最小缓冲
  "channels": 2,
  "input1_gain": 1.0,
  "input2_gain": 1.0,
  "master_gain": 1.0
}
```

### 注意事项
- 可能导致音频断续
- 需要更多CPU资源
- 关闭其他程序

### 适用场景
- 实时演奏
- 现场表演
- 语音监听

## 示例9: 编程接口使用

在自己的Go程序中使用Audio Mixer。

### 代码示例
```go
package main

import (
    "github.com/entropy/audio-mixer/internal/audio"
)

func main() {
    // 初始化
    dm := audio.NewDeviceManager()
    dm.Initialize()
    defer dm.Terminate()

    // 获取设备
    input1, _ := dm.GetDefaultInputDevice()
    output, _ := dm.GetDefaultOutputDevice()

    // 创建混音器
    config := audio.DefaultMixerConfig()
    config.Input1Device = input1
    config.OutputDevice = output
    config.Input1Gain = 1.0

    mixer, _ := audio.NewMixer(config)

    // 启动
    mixer.Start()
    defer mixer.Stop()

    // 动态调节音量
    mixer.SetInput1Gain(0.8)

    // 读取电平
    level := mixer.GetInput1Level()
    latency := mixer.GetLatency()

    // ... 你的逻辑 ...
}
```

## 常见配置参数

### 音量增益参考

| 场景 | Input1 (Mic) | Input2 (App) | Master |
|------|-------------|--------------|---------|
| 语音+音乐 | 1.0 | 0.3 | 1.0 |
| 游戏直播 | 1.2 | 0.8 | 0.9 |
| 播客录制 | 1.0 | 0.15 | 1.0 |
| 在线教学 | 1.0 | 0.5 | 1.0 |
| 音乐DJ | 0.8 | 0.8 | 1.0 |

### 缓冲区大小参考

| Buffer Size | 延迟 | 稳定性 | 适用场景 |
|-------------|------|--------|----------|
| 128 | ~3ms | 低 | 实时演奏 |
| 256 | ~5ms | 中 | 低延迟通话 |
| 512 | ~11ms | 高 | 通用场景 ⭐ |
| 1024 | ~21ms | 很高 | 录音/直播 |
| 2048 | ~43ms | 最高 | 批处理 |

### 采样率参考

| Sample Rate | 音质 | CPU | 适用场景 |
|-------------|------|-----|----------|
| 44100 Hz | 标准CD | 低 | 通用 |
| 48000 Hz | 专业音频 | 中 | 推荐 ⭐ |
| 96000 Hz | 高保真 | 高 | 音乐制作 |
| 192000 Hz | 超高保真 | 很高 | 母带处理 |

## 故障排除示例

### 问题: 听到自己的声音有延迟

**原因**: 缓冲区太大或者有回声

**解决方案**:
```json
{
  "buffer_size": 256,
  "master_gain": 0.8
}
```

并确保:
- 关闭操作系统的"监听此设备"功能
- 使用耳机避免物理反馈

### 问题: 音频断断续续

**原因**: CPU过载或缓冲区太小

**解决方案**:
```json
{
  "buffer_size": 1024
}
```

并:
- 关闭其他占用CPU的程序
- 检查系统资源

### 问题: 音乐音量太大

**解决方案**:
```json
{
  "input2_gain": 0.2
}
```

或在运行时动态调整:
```bash
# 重启audio-mixer并输入新的gain值
```

## 高级用法

### 使用环境变量

```bash
# 启用PortAudio调试
export PA_DEBUG=1
./audio-mixer

# 指定配置文件
export AUDIO_MIXER_CONFIG=/path/to/config.json
./audio-mixer
```

### 使用不同配置文件

```bash
# 备份当前配置
cp ~/.audio-mixer/config.json ~/.audio-mixer/config-gaming.json

# 使用不同配置
cp ~/.audio-mixer/config-podcast.json ~/.audio-mixer/config.json
./audio-mixer
```

### 脚本化启动

```bash
#!/bin/bash
# start-mixer.sh

# 检查依赖
command -v audio-mixer >/dev/null 2>&1 || {
    echo "Audio Mixer not found!"
    exit 1
}

# 设置配置
cat > ~/.audio-mixer/config.json << EOF
{
  "input1_device_index": 0,
  "input2_device_index": 5,
  "output_device_index": 6,
  "input1_gain": 1.0,
  "input2_gain": 0.3,
  "master_gain": 1.0
}
EOF

# 启动
./audio-mixer
```

## 参考资源

- [README.md](README.md) - 基本使用指南
- [INSTALL.md](INSTALL.md) - 安装说明
- [ARCHITECTURE.md](ARCHITECTURE.md) - 架构设计
- [PortAudio文档](http://www.portaudio.com/docs/)
- [BlackHole虚拟设备](https://github.com/ExistentialAudio/BlackHole)

## 贡献你的示例

如果你有好的使用示例,欢迎提交PR添加到本文档!

格式:
```markdown
## 示例X: 标题

### 场景描述

### 配置

### 使用技巧
```
