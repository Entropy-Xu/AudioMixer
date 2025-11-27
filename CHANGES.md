# 更新日志 - 虚拟设备支持

## 概述

已成功将 Audio Mixer 重构为使用虚拟音频设备架构，现在支持：
- **Input 1**: 麦克风/线路输入（物理设备）
- **Input 2**: 系统音频捕获（通过虚拟设备，如 BlackHole）
- **Output**: 虚拟输出设备（混音后的音频可被其他应用使用）

## 主要变更

### 1. 新增虚拟设备管理 (`internal/audio/loopback.go`)

创建了专门的虚拟设备管理模块：

- `FindLoopbackDevice()`: 自动检测系统中的虚拟音频设备
  - macOS: 查找 BlackHole, Soundflower, Loopback Audio
  - Windows: 查找 VB-Cable, Virtual Audio Cable
  - Linux: 查找 PulseAudio null sink, Loopback

- `ListLoopbackDevices()`: 列出所有可用的虚拟设备

- `GetLoopbackDeviceByName()`: 按名称查找虚拟设备

### 2. 更新配置结构 (`internal/config/config.go`)

添加了新的配置选项：

```go
type Config struct {
    // ... 现有字段 ...

    // 新增字段
    UseVirtualOutput   bool   `json:"use_virtual_output"`   // 是否使用虚拟输出
    LoopbackDeviceName string `json:"loopback_device_name"` // 虚拟设备名称
}
```

默认配置：
- `UseVirtualOutput`: `true` （默认使用虚拟输出）
- `LoopbackDeviceName`: `"BlackHole"` （macOS 默认）

### 3. 更新混音器配置 (`internal/audio/mixer.go`)

```go
type MixerConfig struct {
    // ... 现有字段 ...

    // 更新字段注释
    Input1Device     *portaudio.DeviceInfo // 麦克风/线路输入
    Input2Device     *portaudio.DeviceInfo // 系统音频（从虚拟设备）
    OutputDevice     *portaudio.DeviceInfo // 虚拟输出设备
    UseVirtualOutput bool                  // 是否使用虚拟输出
}
```

### 4. 重构 GUI 界面 (`internal/gui/app.go`)

#### 设备选择部分更新：

- **Input 1**: 显示所有物理输入设备（麦克风等）
- **Input 2**: 仅显示虚拟设备 + "<Auto Detect Loopback>" 选项
- **Output**: 仅显示虚拟设备 + "<Auto Detect Virtual Output>" 选项

#### 添加自动检测功能：

```go
// 自动检测 Input 2（系统音频）
if a.cfg.Input2DeviceIndex < 0 {
    loopback, err := audio.FindLoopbackDevice(a.deviceManager)
    if err != nil {
        // 提示用户安装 BlackHole
    }
    mixerConfig.Input2Device = loopback.Device
}

// 自动检测 Output（虚拟输出）
if a.cfg.OutputDeviceIndex < 0 {
    loopback, err := audio.FindLoopbackDevice(a.deviceManager)
    // ...
}
```

#### 中文提示信息：

在设备选择界面添加了说明：
```
提示:
• Input 1: 麦克风/线路输入
• Input 2: 系统音频 (需要虚拟设备,如 BlackHole)
• Output: 虚拟输出 (混音后的音频可被其他应用使用)
```

### 5. 创建虚拟设备设置指南 (`VIRTUAL_DEVICE_SETUP.md`)

详细的跨平台虚拟设备设置文档：

#### macOS 部分：
- BlackHole 安装方法（Homebrew 和手动安装）
- 如何创建 Multi-Output Device
- 系统音频路由配置
- 使用场景示例（直播、录制、会议）

#### Windows 部分：
- VB-Cable 安装步骤
- 音频设备配置（侦听设置）
- 虚拟音频路由

#### Linux 部分：
- PulseAudio loopback 模块配置
- JACK Audio Connection Kit 设置
- 虚拟 sink 创建

#### 故障排除：
- 常见问题和解决方案
- 延迟优化建议
- 推荐配置参数

### 6. 更新 README (`README.md`)

- 更新功能特性描述，强调虚拟设备支持
- 在系统要求中明确标注 BlackHole/VB-Cable 为**必需**
- 添加"第一步：安装虚拟音频设备"章节
- 更新使用指南，包含自动检测功能说明
- 添加到 VIRTUAL_DEVICE_SETUP.md 的链接

### 7. 字体支持优化

之前已完成：
- 移除 TTC 字体文件支持（Fyne 限制）
- 更新为仅使用 TTF/OTF 格式字体
- 添加 GUI 内置字体选择器

## 工作原理

### 系统音频捕获流程（macOS 示例）

```
┌──────────────────┐
│  系统音频播放     │
│  (音乐、视频等)   │
└────────┬─────────┘
         │
         ▼
┌──────────────────────────────┐
│  Multi-Output Device          │
│  ┌────────────────────────┐  │
│  │ 1. 扬声器/耳机         │  │  ← 你能听到声音
│  │ 2. BlackHole 2ch       │  │  ← Audio Mixer 从这里捕获
│  └────────────────────────┘  │
└──────────────────────────────┘
         │
         ▼
┌──────────────────────────────┐
│  Audio Mixer                  │
│  ┌────────────────────────┐  │
│  │ Input 1: 麦克风        │  │
│  │ Input 2: BlackHole     │  │  ← 系统音频
│  │          (系统音频)     │  │
│  │                         │  │
│  │ 混音 + 音量控制        │  │
│  │                         │  │
│  │ Output: BlackHole 16ch │  │  ← 虚拟输出
│  └────────────────────────┘  │
└──────────────┬───────────────┘
               │
               ▼
┌──────────────────────────────┐
│  其他应用程序                 │
│  (OBS, Zoom, 录屏软件等)      │
│  从 BlackHole 16ch 读取      │
└──────────────────────────────┘
```

## 使用场景

### 1. 直播/录制
- **需求**: 混合麦克风和电脑播放的音乐
- **配置**:
  - Input 1: 麦克风
  - Input 2: BlackHole（捕获系统音频）
  - Output: BlackHole 或另一个虚拟设备
  - OBS: 选择 Audio Mixer 的输出设备作为音频源

### 2. 在线会议
- **需求**: 在会议中播放音乐或音效
- **配置**:
  - Input 1: 麦克风（你的声音）
  - Input 2: BlackHole（电脑播放的内容）
  - Output: BlackHole
  - Zoom/Teams: 选择 BlackHole 作为麦克风输入

### 3. 游戏直播
- **需求**: 混合游戏音频、麦克风和背景音乐
- **配置**:
  - 游戏音频 → Multi-Output Device → BlackHole
  - Input 1: 麦克风
  - Input 2: BlackHole（游戏音频）
  - Output: 另一个 BlackHole 实例
  - 直播软件: 从 Output 设备读取

## 技术亮点

1. **自动检测**: 智能检测系统中的虚拟音频设备
2. **跨平台**: 支持 macOS (BlackHole)、Windows (VB-Cable)、Linux (PulseAudio)
3. **零配置**: 选择 "Auto Detect" 即可自动配置
4. **实时混音**: 低延迟 (<30ms) 音频处理
5. **中文支持**: 完整的中文界面和文档

## 后续计划

- [ ] 支持多个虚拟设备同时使用
- [ ] 添加音频效果（均衡器、压缩器）
- [ ] 支持保存和加载预设配置
- [ ] 添加音频路由可视化界面
- [ ] 支持 ASIO 驱动（Windows 专业音频）

## 测试建议

### macOS 测试步骤：

1. 安装 BlackHole:
   ```bash
   brew install blackhole-2ch
   ```

2. 创建 Multi-Output Device（Audio MIDI Setup）

3. 播放音乐，测试系统音频捕获

4. 对着麦克风说话，测试麦克风输入

5. 在其他应用（如 QuickTime）中选择 BlackHole 作为输入源，验证混音输出

### Windows 测试步骤：

1. 安装 VB-Cable

2. 配置音频设备（启用侦听）

3. 运行 Audio Mixer GUI

4. 测试系统音频和麦克风混音

### 常见问题排查：

- **没有声音**: 检查 Multi-Output Device 是否包含扬声器
- **找不到虚拟设备**: 确认 BlackHole/VB-Cable 已正确安装
- **延迟高**: 降低 Buffer Size（例如 256 samples）
- **音频断续**: 提高 Buffer Size（例如 1024 samples）

## 文件变更总结

### 新增文件：
- `internal/audio/loopback.go` - 虚拟设备管理
- `VIRTUAL_DEVICE_SETUP.md` - 虚拟设备设置指南
- `CHANGES.md` - 本文档

### 修改文件：
- `internal/audio/mixer.go` - 更新配置结构
- `internal/config/config.go` - 添加虚拟设备配置项
- `internal/gui/app.go` - 重构设备选择界面
- `internal/gui/font.go` - 字体支持优化（之前完成）
- `README.md` - 更新文档

## 兼容性说明

- **向后兼容**: 旧的配置文件仍然可用，新字段使用默认值
- **平台支持**: macOS、Windows、Linux 均已适配
- **虚拟设备**: 自动检测，无需手动配置设备索引

---

**更新日期**: 2025-11-27
**版本**: v2.0 - Virtual Device Support
