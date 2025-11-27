# 虚拟设备设置指南 (Virtual Device Setup Guide)

本文档将帮助你设置虚拟音频设备，以便 Audio Mixer 可以：
1. 捕获系统音频（其他应用播放的声音）作为 Input 2
2. 将混音后的音频输出到虚拟设备，供其他应用使用

## macOS 设置

### 1. 安装 BlackHole

BlackHole 是一个免费、开源的虚拟音频驱动，用于在应用之间路由音频。

#### 使用 Homebrew 安装（推荐）

```bash
# 安装 2 通道版本（立体声）
brew install blackhole-2ch

# 或者安装 16 通道版本（如果需要多声道支持）
brew install blackhole-16ch
```

#### 手动安装

1. 访问 [BlackHole GitHub Releases](https://github.com/ExistentialAudio/BlackHole/releases)
2. 下载 `BlackHole2ch.vX.X.X.pkg`（立体声版本）
3. 运行安装包并按照提示安装
4. 重启电脑（推荐）

### 2. 创建 Multi-Output Device（用于捕获系统音频）

为了同时听到系统音频并让 Audio Mixer 捕获它，需要创建一个 Multi-Output Device：

1. 打开 **Audio MIDI Setup**（音频 MIDI 设置）：
   - 在 Spotlight 中搜索 "Audio MIDI Setup"
   - 或者在 `/Applications/Utilities/` 中找到它

2. 点击左下角的 **+** 按钮，选择 **Create Multi-Output Device**

3. 在 Multi-Output Device 中勾选：
   - ✅ **你的扬声器/耳机**（例如：MacBook Pro Speakers 或 External Headphones）
   - ✅ **BlackHole 2ch**

4. 将这个 Multi-Output Device 命名为 "System Audio + BlackHole"

5. **重要**: 在系统设置中设置音频输出：
   - 打开 **System Settings** → **Sound** → **Output**
   - 选择 "**System Audio + BlackHole**" 作为输出设备
   - 这样系统音频会同时输出到扬声器和 BlackHole

### 3. 配置 Audio Mixer

现在你的设备配置应该是：

- **Input 1 (麦克风)**: 选择你的麦克风设备
- **Input 2 (系统音频)**: 选择 "**BlackHole 2ch**"（从 Multi-Output Device 捕获）
- **Output (虚拟输出)**: 选择 "**BlackHole 2ch**"（或另一个 BlackHole 实例）

### 4. 使用场景示例

#### 场景 1: 直播/录制（麦克风 + 系统音频）

1. 系统音频输出: **System Audio + BlackHole** (Multi-Output Device)
2. Audio Mixer 设置:
   - Input 1: 你的麦克风
   - Input 2: BlackHole 2ch（捕获系统音频）
   - Output: BlackHole 2ch 或另一个虚拟设备
3. OBS/录屏软件: 选择 Audio Mixer 的输出设备作为音频源

#### 场景 2: 在线会议（混合音频源）

1. 系统音频输出: **System Audio + BlackHole**
2. Audio Mixer 设置:
   - Input 1: 麦克风（你的声音）
   - Input 2: BlackHole 2ch（电脑播放的音乐/音效）
   - Output: BlackHole 16ch（如果安装了）
3. Zoom/Teams: 选择 BlackHole 16ch 作为麦克风输入

---

## Windows 设置

### 1. 安装 VB-Cable 或 Virtual Audio Cable

#### VB-Cable (免费)

1. 下载 [VB-Audio Virtual Cable](https://vb-audio.com/Cable/)
2. 解压下载的文件
3. 右键点击 `VBCABLE_Setup_x64.exe` → **以管理员身份运行**
4. 安装完成后**重启电脑**

#### Virtual Audio Cable (付费，功能更强)

1. 购买并下载 [Virtual Audio Cable](https://vac.muzychenko.net/en/)
2. 运行安装程序
3. 重启电脑

### 2. 配置音频设备

1. 右键点击任务栏的音量图标 → **声音设置**
2. 在 **输出** 中选择你的扬声器
3. 打开 **声音控制面板** (Sound Control Panel)：
   - 在 **播放** 标签中，找到 "**CABLE Input**"
   - 右键 → **属性** → **侦听** 标签
   - 勾选 "**侦听此设备**"
   - 在下拉菜单中选择你的**真实扬声器**
   - 这样可以同时输出到 VB-Cable 和扬声器

### 3. 配置 Audio Mixer

- **Input 1**: 你的麦克风
- **Input 2**: CABLE Output（捕获系统音频）
- **Output**: CABLE Input（虚拟输出）

---

## Linux 设置

### 1. 使用 PulseAudio Loopback Module

```bash
# 加载 loopback 模块（将一个设备的输出路由到另一个设备的输入）
pactl load-module module-loopback latency_msec=1

# 列出所有音频设备
pactl list sources
pactl list sinks

# 创建一个虚拟 sink（null sink）
pactl load-module module-null-sink sink_name=VirtualOutput sink_properties=device.description=VirtualOutput
```

### 2. 使用 JACK Audio Connection Kit

JACK 提供了更专业的音频路由功能：

```bash
# Ubuntu/Debian
sudo apt install jackd2 qjackctl

# Fedora
sudo dnf install jack-audio-connection-kit qjackctl

# 启动 JACK
qjackctl
```

### 3. 配置 Audio Mixer

- **Input 1**: 你的麦克风
- **Input 2**: monitor of null sink（监控虚拟 sink）
- **Output**: VirtualOutput（虚拟 sink）

---

## 故障排除

### macOS

**问题**: 安装 BlackHole 后没有声音

**解决方案**:
1. 确保创建了 Multi-Output Device 并勾选了你的扬声器
2. 在系统设置中选择 Multi-Output Device 作为输出
3. 检查扬声器音量没有被静音

**问题**: Audio Mixer 提示 "未找到虚拟设备"

**解决方案**:
1. 确认 BlackHole 已正确安装：在 Audio MIDI Setup 中应该能看到 BlackHole 2ch
2. 重启 Audio Mixer 应用
3. 如果仍然不行，重启电脑

### Windows

**问题**: 安装 VB-Cable 后无法播放音频

**解决方案**:
1. 在声音控制面板中检查 "CABLE Input" 的 "侦听" 设置
2. 确保 "侦听此设备" 已勾选，并选择了正确的扬声器

**问题**: 延迟太高

**解决方案**:
1. 在 Audio Mixer 中降低 Buffer Size（例如从 512 降到 256）
2. 在 VB-Cable 控制面板中调整缓冲区大小
3. 关闭不必要的音频应用

### Linux

**问题**: PulseAudio loopback 延迟高

**解决方案**:
```bash
# 卸载现有 loopback 模块
pactl unload-module module-loopback

# 重新加载，使用更低的延迟
pactl load-module module-loopback latency_msec=1
```

---

## 推荐配置

### 直播/录制场景
- **Sample Rate**: 48000 Hz
- **Buffer Size**: 512 samples（延迟 ~10ms）
- **Channels**: 2 (Stereo)

### 音乐制作场景
- **Sample Rate**: 48000 或 96000 Hz
- **Buffer Size**: 256 samples（延迟 ~5ms）
- **Channels**: 2 (Stereo)

### 在线会议场景
- **Sample Rate**: 44100 或 48000 Hz
- **Buffer Size**: 1024 samples（延迟 ~20ms，更稳定）
- **Channels**: 2 (Stereo) 或 1 (Mono，节省带宽）

---

## 参考链接

- [BlackHole GitHub](https://github.com/ExistentialAudio/BlackHole)
- [VB-Audio Virtual Cable](https://vb-audio.com/Cable/)
- [Virtual Audio Cable](https://vac.muzychenko.net/en/)
- [PulseAudio Documentation](https://www.freedesktop.org/wiki/Software/PulseAudio/Documentation/)
- [JACK Audio Connection Kit](https://jackaudio.org/)
