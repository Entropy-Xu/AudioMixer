# 更新总结 - 自定义虚拟输出设备

## 更新内容

根据你的需求 "把虚拟输出调整为自己设置一个名字的虚拟设备来进行混合后的输出"，已完成以下更新：

### 主要变更

#### 1. **移除下拉列表，改为文本输入**

**之前:**
```
Output: [下拉列表选择]
  ├─ <Auto Detect Virtual Output>
  ├─ [Virtual] BlackHole 2ch
  └─ [Virtual] BlackHole 16ch
```

**现在:**
```
Output (虚拟输出设备名称):
  ┌────────────────────────────────┐
  │ BlackHole 2ch                  │  [检测设备]
  └────────────────────────────────┘
```

#### 2. **新增设备检测功能**

点击 "检测设备 (Detect)" 按钮可以：
- 验证输入的设备名称是否存在
- 显示检测结果：
  - ✓ 找到设备: [设备名]
  - ✗ 未找到设备: [错误信息]

#### 3. **灵活的设备名称匹配**

支持部分匹配（不区分大小写）：
- 输入 `BlackHole` → 匹配 `BlackHole 2ch`, `BlackHole 16ch`
- 输入 `BlackHole 2ch` → 精确匹配 `BlackHole 2ch`
- 输入 `VB` → 匹配 `VB-Cable`, `VB-Audio`

#### 4. **配置持久化**

设备名称保存在配置文件中：
```json
{
  "loopback_device_name": "BlackHole 2ch",
  "use_virtual_output": true
}
```

## 使用示例

### 基本用法

1. **启动 GUI**:
   ```bash
   go build -o audio-mixer-gui ./cmd/gui
   ./audio-mixer-gui
   ```

2. **配置输出设备**:
   - 在 "Output (虚拟输出设备名称)" 输入框中输入: `BlackHole 2ch`
   - 点击 "检测设备 (Detect)" 验证
   - 看到 "✓ 找到设备: BlackHole 2ch" 表示成功

3. **开始混音**:
   - 点击 "Start Mixer"
   - 混音后的音频将输出到 BlackHole 2ch

### 高级场景

#### 场景 1: 多实例混音（链式）

**实例 1 - 基础混音:**
```
Input 1: 麦克风
Input 2: 系统音频 (通过 BlackHole 2ch)
Output: BlackHole 16ch  ← 输出到 16 通道设备
```

**实例 2 - 二次混音:**
```
Input 1: BlackHole 16ch  ← 从第一个混音器读取
Input 2: 背景音乐
Output: BlackHole 64ch   ← 输出到 64 通道设备
```

**最终应用 (OBS):**
```
音频源: BlackHole 64ch  ← 使用最终混音
```

#### 场景 2: 多个虚拟设备同时使用

如果你安装了多个虚拟设备：
- BlackHole 2ch (用于系统音频捕获)
- BlackHole 16ch (用于混音输出)

配置：
```
Input 2: BlackHole 2ch
Output: BlackHole 16ch
```

#### 场景 3: Windows VB-Cable

```
Output 设备名称: CABLE Input
或
Output 设备名称: VB-Cable
```

## 技术细节

### 代码变更

**文件: `internal/gui/app.go`**

1. **添加输入框:**
```go
a.outputNameEntry = widget.NewEntry()
a.outputNameEntry.SetPlaceHolder("输入虚拟设备名称，如: BlackHole 2ch")
a.outputNameEntry.OnChanged = func(value string) {
    a.cfg.LoopbackDeviceName = value
}
```

2. **添加检测按钮:**
```go
detectButton := widget.NewButton("检测设备 (Detect)", func() {
    loopback, err := audio.GetLoopbackDeviceByName(a.deviceManager, a.outputNameEntry.Text)
    if err != nil {
        a.statusLabel.SetText(fmt.Sprintf("未找到设备: %v", err))
    } else {
        a.statusLabel.SetText(fmt.Sprintf("✓ 找到设备: %s", loopback.Name))
    }
})
```

3. **更新 startMixer:**
```go
if a.cfg.LoopbackDeviceName != "" {
    loopback, err := audio.GetLoopbackDeviceByName(a.deviceManager, a.cfg.LoopbackDeviceName)
    if err != nil {
        a.statusLabel.SetText(fmt.Sprintf("未找到设备 '%s': %v", a.cfg.LoopbackDeviceName, err))
        return
    }
    mixerConfig.OutputDevice = loopback.Device
}
```

### 已修复问题

✅ **通道数错误**: 修复了 "Invalid number of channels" 错误
- 现在会根据设备的实际支持通道数动态调整
- 单声道设备使用 1 通道，立体声设备使用 2 通道

✅ **设备索引混淆**: 使用设备名称而非索引
- 配置文件更易读
- 跨设备配置更容易迁移

✅ **灵活性**: 支持任意虚拟设备
- 不限于自动检测到的设备列表
- 支持自定义虚拟设备名称

## 界面截图描述

```
┌─────────────────────────────────────────────┐
│ Audio Mixer                                  │
├─────────────────────────────────────────────┤
│ 设备配置 (Devices)                           │
│                                              │
│ Input 1 (麦克风):   [MacBook Pro Microphone]│
│ Input 2 (系统音频): [<Auto Detect Loopback>]│
│                                              │
│ Output (虚拟输出设备名称):                   │
│ ┌──────────────────────────────────────────┐│
│ │ BlackHole 2ch                   [检测设备]││
│ └──────────────────────────────────────────┘│
│                                              │
│ 提示:                                        │
│ • Input 1: 麦克风/线路输入                   │
│ • Input 2: 系统音频 (需要虚拟设备)           │
│ • Output: 输入虚拟设备名称                   │
│ • 常见设备名: BlackHole 2ch, VB-Cable        │
│                                              │
│ ... (音量控制、电平表等)                     │
│                                              │
│ [Start Mixer]  [Stop Mixer]                 │
│                                              │
│ Status: ✓ 找到设备: BlackHole 2ch           │
└─────────────────────────────────────────────┘
```

## 文档

新增和更新的文档：

1. **[CUSTOM_OUTPUT_DEVICE.md](CUSTOM_OUTPUT_DEVICE.md)** - 自定义输出设备完整指南
2. **[VIRTUAL_DEVICE_SETUP.md](VIRTUAL_DEVICE_SETUP.md)** - 虚拟设备安装配置
3. **[README.md](README.md)** - 更新使用说明
4. **[UPDATE_SUMMARY.md](UPDATE_SUMMARY.md)** - 本文档

## 测试建议

### macOS 测试步骤:

1. 确保已安装 BlackHole:
   ```bash
   brew install blackhole-2ch
   # 或
   brew install blackhole-16ch
   ```

2. 编译并运行:
   ```bash
   go build -o audio-mixer-gui ./cmd/gui
   ./audio-mixer-gui
   ```

3. 测试输入:
   - 输入: `BlackHole 2ch`
   - 点击 "检测设备"
   - 应显示: "✓ 找到设备: BlackHole 2ch"

4. 测试部分匹配:
   - 输入: `BlackHole`
   - 点击 "检测设备"
   - 应找到第一个匹配的设备 (2ch 或 16ch)

5. 测试错误情况:
   - 输入: `NonExistentDevice`
   - 点击 "检测设备"
   - 应显示: "未找到设备: loopback device 'NonExistentDevice' not found"

6. 测试混音:
   - 配置 Input 1 和 Input 2
   - 设置 Output 为 `BlackHole 2ch`
   - 点击 "Start Mixer"
   - 使用 QuickTime 或其他软件从 BlackHole 2ch 录音验证

### Windows 测试步骤:

1. 安装 VB-Cable

2. 运行程序并输入: `CABLE Input` 或 `VB-Cable`

3. 验证设备检测和混音功能

## 配置文件示例

**位置:** `~/.audio-mixer/config.json`

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

## 优势总结

与之前的实现相比：

| 特性 | 之前 (下拉列表) | 现在 (自定义名称) |
|------|----------------|------------------|
| 灵活性 | 仅限检测到的设备 | 任意虚拟设备 |
| 易用性 | 需要滚动查找 | 直接输入名称 |
| 配置可读性 | 使用设备索引 | 使用设备名称 |
| 跨设备兼容 | 索引可能不同 | 名称通用 |
| 设备验证 | 无 | 有检测按钮 |
| 多设备支持 | 有限 | 完全支持 |

## 后续改进计划

- [ ] 添加设备名称自动补全
- [ ] 显示最近使用的设备列表
- [ ] 支持从系统设备列表中选择（保留当前输入框）
- [ ] 添加设备属性显示（通道数、采样率等）
- [ ] 命令行参数支持: `--output-device "BlackHole 2ch"`

---

**更新时间:** 2025-11-27
**版本:** v2.1 - Custom Output Device Name
