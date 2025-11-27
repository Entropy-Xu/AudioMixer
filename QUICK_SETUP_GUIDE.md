# 快速配置指南 - 捕获特定应用音频

## 目标

实现类似 OBS 的功能：
- ✅ 捕获特定应用（如音乐播放器）的音频
- ✅ 混合麦克风和应用音频
- ✅ 输出到独立的虚拟设备供其他软件使用

## macOS 快速配置（5分钟）

### 第1步：安装多个 BlackHole 实例

```bash
# 安装 2ch 版本（用于捕获应用音频）
brew install blackhole-2ch

# 安装 16ch 版本（用于混音输出）
brew install blackhole-16ch
```

安装后你会有：
- **BlackHole 2ch** ← 应用音频路由到这里
- **BlackHole 16ch** ← Audio Mixer 输出到这里

### 第2步：创建 Multi-Output Device

1. 打开 **Audio MIDI Setup** (音频 MIDI 设置)
   - Spotlight 搜索: `Audio MIDI Setup`
   - 或路径: `/Applications/Utilities/Audio MIDI Setup.app`

2. 点击左下角 **+** 按钮

3. 选择 **Create Multi-Output Device**

4. 配置 Multi-Output Device：
   - ✅ 勾选 **BlackHole 2ch**
   - ✅ 勾选你的**扬声器/耳机**（这样你可以听到声音）
   - 重命名为: `Music + Speakers`

5. **重要设置**：
   - 右键点击你的扬声器 → 设为 **Clock Source**
   - 这样可以避免音频漂移

### 第3步：配置音乐播放器

在你想捕获的应用中设置音频输出：

**Spotify:**
1. Preferences → Playback
2. Audio Output → `Music + Speakers`

**Apple Music / iTunes:**
1. 点击音量图标旁边的设备选择
2. 选择 `Music + Speakers`

**VLC / QuickTime:**
1. Audio → Audio Device
2. 选择 `Music + Speakers`

**系统默认（所有应用）:**
1. System Settings → Sound → Output
2. 选择 `Music + Speakers`

### 第4步：配置 Audio Mixer

1. 运行 Audio Mixer GUI:
   ```bash
   ./audio-mixer-gui
   ```

2. 配置设备：
   ```
   Input 1 (麦克风): MacBook Pro Microphone
   Input 2 (系统音频): <Auto Detect Loopback>  (会自动找到 BlackHole 2ch)
   Output (虚拟输出): BlackHole 16ch
   ```

3. 调节音量：
   - Input 1 (麦克风): 100%
   - Input 2 (音乐): 50-80% (根据需要调整)
   - Master: 100%

4. 点击 **Start Mixer**

### 第5步：在其他应用中使用混音

**OBS Studio:**
1. Sources → Add → Audio Input Capture
2. Device → `BlackHole 16ch`
3. 现在 OBS 会录制麦克风+音乐的混音

**Zoom / Teams:**
1. Audio Settings
2. Microphone → `BlackHole 16ch`
3. 你的声音+音乐会被传输

**录音软件 (Audacity / Logic Pro):**
1. Input Device → `BlackHole 16ch`
2. 录制即可

## 完整音频流图

```
┌─────────────────────────────┐
│ Spotify / Apple Music       │
│ 输出: Music + Speakers      │
└──────────┬──────────────────┘
           │
           ├────────────────────────┐
           │                        │
           ▼                        ▼
   ┌───────────────┐     ┌──────────────────┐
   │ BlackHole 2ch │     │ 你的扬声器/耳机  │ ◄─ 你能听到音乐
   └───────┬───────┘     └──────────────────┘
           │
           │ (捕获音乐)
           │
           ▼
   ┌────────────────────────────┐
   │ Audio Mixer                │
   │ ┌────────────────────────┐ │
   │ │ Input 1: 麦克风        │ │
   │ │ Input 2: BlackHole 2ch │ │◄─ 音乐在这里
   │ │                         │ │
   │ │ 混音 (麦克风 + 音乐)   │ │
   │ │                         │ │
   │ │ Output: BlackHole 16ch │ │
   │ └────────────────────────┘ │
   └────────────┬───────────────┘
                │
                │ (混音后的音频)
                │
                ▼
   ┌────────────────────────────┐
   │ OBS / Zoom / 录音软件       │
   │ 输入: BlackHole 16ch        │
   └────────────────────────────┘
```

## 实际使用场景

### 场景 1: 直播游戏 + 背景音乐

1. **游戏音频**: 系统默认输出
2. **音乐播放器**: 输出到 `Music + Speakers`
3. **Audio Mixer**:
   - Input 1: 麦克风
   - Input 2: BlackHole 2ch (音乐)
   - Output: BlackHole 16ch
4. **OBS**: 音频源 BlackHole 16ch

结果：直播中有麦克风、游戏声音、背景音乐

### 场景 2: 播客录制 + 片头音乐

1. **音乐播放器**: 输出到 `Music + Speakers`
2. **Audio Mixer**:
   - Input 1: 麦克风
   - Input 2: BlackHole 2ch (音乐)
   - 音乐音量: 30%（背景音乐）
   - 麦克风音量: 100%
3. **Audacity**: 录音源 BlackHole 16ch

### 场景 3: 在线会议 + 演示音频

1. **演示视频播放器**: 输出到 `Music + Speakers`
2. **Audio Mixer**: 混合麦克风和演示音频
3. **Zoom**: 麦克风选择 BlackHole 16ch

## 高级技巧

### 技巧 1: 分离不同应用的音频

如果你想分别控制多个应用的音量：

1. 安装更多 BlackHole 实例:
   ```bash
   brew install blackhole-64ch
   ```

2. 创建多个 Multi-Output Device:
   - `Music + Speakers` (BlackHole 2ch + 扬声器)
   - `Game + Speakers` (BlackHole 16ch + 扬声器)
   - `Discord + Speakers` (BlackHole 64ch + 扬声器)

3. 使用多个 Audio Mixer 实例混合

### 技巧 2: 动态切换音频源

创建多个配置文件：

**~/.audio-mixer/config-music.json**
```json
{
  "loopback_device_name": "BlackHole 2ch",
  "input2_device_index": -1
}
```

**~/.audio-mixer/config-game.json**
```json
{
  "loopback_device_name": "BlackHole 16ch",
  "input2_device_index": -1
}
```

切换配置：
```bash
cp ~/.audio-mixer/config-music.json ~/.audio-mixer/config.json
```

### 技巧 3: 使用 AppleScript 自动化

自动设置音乐播放器输出：

```applescript
tell application "Music"
    set sound output to "Music + Speakers"
end tell
```

## 故障排除

### ❌ 问题：听不到音乐

**检查项：**
1. Multi-Output Device 是否勾选了扬声器？
2. 扬声器是否设为 Clock Source？
3. 扬声器音量是否静音？

**解决：**
- 重新创建 Multi-Output Device
- 确保同时勾选 BlackHole 2ch 和扬声器

### ❌ 问题：OBS 录制没声音

**检查项：**
1. Audio Mixer 是否正在运行？
2. OBS 音频源是否选择了 BlackHole 16ch？
3. Input 2 是否正确配置为 BlackHole 2ch？

**解决：**
- 在 Audio Mixer 中查看 Input 2 的电平表是否有波动
- 点击 "检测设备" 验证 BlackHole 16ch 是否存在

### ❌ 问题：音频有延迟/不同步

**解决：**
1. 降低 Audio Mixer 的 Buffer Size:
   - 512 → 256 samples
2. 确保 Multi-Output Device 的 Clock Source 设置正确
3. 关闭不必要的音频应用

### ❌ 问题：音质下降/有杂音

**解决：**
1. 提高 Buffer Size: 256 → 512 或 1024
2. 检查采样率是否一致（都使用 48000 Hz）
3. 降低混音器的音量增益，避免削波

## Windows 版本（类似配置）

### 使用 VB-Cable 实现

1. 安装 VB-Cable: https://vb-audio.com/Cable/

2. 配置音乐播放器:
   - 输出设备: CABLE Input

3. Audio Mixer 配置:
   - Input 1: 麦克风
   - Input 2: CABLE Output
   - Output: 另一个虚拟设备

## Linux 版本（PulseAudio）

```bash
# 创建虚拟 sink
pactl load-module module-null-sink sink_name=music sink_properties=device.description=MusicSink

# 创建 loopback
pactl load-module module-loopback source=music.monitor sink=@DEFAULT_SINK@

# Audio Mixer 配置
# Input 2: music.monitor
```

## 总结

通过这个配置，你实现了：
- ✅ **捕获特定应用音频** (音乐播放器 → BlackHole 2ch)
- ✅ **独立虚拟输出** (Audio Mixer → BlackHole 16ch)
- ✅ **灵活混音控制** (独立调节每个源的音量)
- ✅ **保持音频监听** (通过 Multi-Output Device 听到音乐)

这个方案：
- 🆓 完全免费
- 🚀 配置简单（5分钟）
- 🔧 灵活可控
- 🎯 满足专业需求

如果需要更高级的功能（如 GUI 选择应用、自动路由等），可以考虑使用商业软件 Loopback ($99)，或者我可以为 Audio Mixer 添加 ScreenCaptureKit 集成来实现应用选择功能。
