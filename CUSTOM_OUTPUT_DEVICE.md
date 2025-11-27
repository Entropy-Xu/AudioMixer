# 自定义虚拟输出设备名称

## 功能说明

现在 Audio Mixer 支持通过**自定义设备名称**来指定虚拟输出设备，而不是从下拉列表中选择。这样可以更灵活地配置不同的虚拟设备。

## 使用方法

### 1. 在 GUI 中设置

1. 启动 Audio Mixer GUI
2. 在 "Output (虚拟输出设备名称)" 输入框中输入你的虚拟设备名称
3. 点击 "检测设备 (Detect)" 按钮验证设备是否存在
4. 如果显示 "✓ 找到设备: [设备名]"，说明配置正确
5. 点击 "Start Mixer" 开始混音

### 2. 常见设备名称

#### macOS:
```
BlackHole 2ch
BlackHole 16ch
Soundflower (2ch)
Soundflower (64ch)
Loopback Audio
```

#### Windows:
```
CABLE Input
VB-Cable
VB-Audio Virtual Cable
```

#### Linux:
```
pulse
null
Monitor of Built-in Audio
```

## 配置示例

### 场景 1: 使用 BlackHole 2ch 用于简单混音

**设置:**
- Input 1: MacBook Pro Microphone
- Input 2: BlackHole 2ch (系统音频)
- Output: `BlackHole 2ch`

**用途:**
- 混音后的音频输出到 BlackHole 2ch
- OBS 或其他录音软件从 BlackHole 2ch 捕获混音音频

### 场景 2: 使用 BlackHole 16ch 用于多声道

**设置:**
- Input 1: 你的麦克风
- Input 2: BlackHole 2ch (系统音频)
- Output: `BlackHole 16ch`

**用途:**
- 混音后输出到 16 通道虚拟设备
- 支持更多声道的专业应用

### 场景 3: 使用不同虚拟设备实现链式混音

**设置 1 - Audio Mixer 实例 1:**
- Input 1: 麦克风
- Input 2: 系统音频 (通过 Multi-Output Device)
- Output: `BlackHole 2ch`

**设置 2 - Audio Mixer 实例 2:**
- Input 1: BlackHole 2ch (上一个混音器的输出)
- Input 2: 另一个音频源
- Output: `BlackHole 16ch`

**设置 3 - 最终应用 (OBS/Zoom):**
- 音频输入: BlackHole 16ch

## 配置文件

设备名称保存在配置文件中：

**位置:** `~/.audio-mixer/config.json`

**示例:**
```json
{
  "sample_rate": 48000,
  "buffer_size": 512,
  "channels": 2,
  "input1_device_index": 0,
  "input2_device_index": -1,
  "output_device_index": -1,
  "use_virtual_output": true,
  "loopback_device_name": "BlackHole 2ch",
  "input1_gain": 1.0,
  "input2_gain": 1.0,
  "master_gain": 1.0,
  "window_width": 800,
  "window_height": 600,
  "start_minimized": false
}
```

**关键字段说明:**
- `loopback_device_name`: 虚拟输出设备的名称
- `use_virtual_output`: 是否使用虚拟输出（应始终为 `true`）
- `output_device_index`: 已弃用，现在使用设备名称

## 设备名称匹配规则

设备名称匹配使用**部分匹配**（不区分大小写）：

**输入:** `BlackHole`
**匹配到:**
- BlackHole 2ch
- BlackHole 16ch
- BlackHole (any)

**输入:** `BlackHole 2ch`
**匹配到:**
- BlackHole 2ch

**输入:** `VB-Cable`
**匹配到:**
- CABLE Input (VB-Audio Virtual Cable)
- VB-Audio Virtual Cable

## 故障排除

### 问题 1: "未找到设备"

**原因:** 输入的设备名称不存在或拼写错误

**解决方案:**
1. 检查虚拟设备是否已安装
2. 在 macOS 中打开 "Audio MIDI Setup" 查看可用设备
3. 在 Windows 中打开 "声音控制面板" 查看设备列表
4. 尝试使用部分名称，如 `BlackHole` 而不是 `BlackHole 2ch`

### 问题 2: 设备名称包含特殊字符

**示例:** 某些设备名称包含括号、数字等

**解决方案:**
- 完整输入设备名称，包括所有特殊字符
- 例如: `Soundflower (2ch)` 需要包含括号

### 问题 3: 有多个同名设备

**解决方案:**
- 使用更具体的名称
- 例如: 使用 `BlackHole 2ch` 而不是 `BlackHole`

### 问题 4: 配置保存后无法加载

**解决方案:**
1. 手动编辑 `~/.audio-mixer/config.json`
2. 确保 `loopback_device_name` 字段格式正确
3. 如果无法修复，删除配置文件让程序重新生成

## 高级用法

### 使用环境变量设置默认设备

你可以通过修改默认配置来设置默认虚拟设备：

**编辑:** `internal/config/config.go`

```go
func DefaultConfig() *Config {
    return &Config{
        // ...
        LoopbackDeviceName: "你的默认设备名", // 修改这里
        // ...
    }
}
```

### 命令行模式指定设备（未来功能）

计划支持：
```bash
./audio-mixer-gui --output-device "BlackHole 16ch"
```

## 优势

相比之前的下拉列表选择方式：

✅ **更灵活**: 可以指定任何虚拟设备，不限于自动检测到的设备
✅ **更精确**: 直接输入完整设备名称，避免索引混淆
✅ **更易配置**: 配置文件中使用设备名称而非索引，更易读
✅ **跨设备兼容**: 配置文件在不同机器上更容易迁移
✅ **支持检测**: 提供"检测设备"按钮，验证设备是否存在

## 相关文档

- [虚拟设备设置指南](VIRTUAL_DEVICE_SETUP.md) - 如何安装和配置虚拟音频设备
- [README](README.md) - 主要使用文档
- [CHANGES](CHANGES.md) - 详细变更日志
