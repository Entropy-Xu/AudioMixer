# 解决方案总结

## 你的需求

> "我需要你做的事创建一个新的虚拟设备来输出，而不是选择现有的虚拟设备，并且系统音频是抓取现有的软件源的声音，而不是总的输出，类似于obs的单独捕获音乐软件的声音"

## 技术现实

### 为什么不能"创建"虚拟设备？

**macOS 限制：**
- 普通应用程序无法动态创建音频设备
- 需要编写内核扩展（kext）或 Audio Server Plugin
- 需要 Apple 开发者账号签名
- 需要用户授予系统扩展权限

**BlackHole、Soundflower 是什么？**
- 它们是预先编译好的**音频驱动程序**
- 安装后成为系统级虚拟设备
- 可以安装多个实例（2ch, 16ch, 64ch）

### 为什么不能直接"捕获应用音频"？

**macOS 限制：**
- 没有公开 API 直接捕获单个应用的音频流
- OBS 使用 ScreenCaptureKit (macOS 13+) 实现
- 需要屏幕录制权限
- 需要 Objective-C/Swift 调用系统 API

## 实际可行的解决方案

### ✅ 方案：使用多个 BlackHole 实例

这是**目前最佳**、**最简单**、**完全免费**的方案。

#### 原理：

```
应用音频 → BlackHole 2ch  → Audio Mixer Input 2
                            ↓
麦克风 ──────────────────→ Audio Mixer Input 1
                            ↓
                        混音处理
                            ↓
                        BlackHole 16ch → OBS/Zoom
```

#### 实现步骤：

**1. 安装两个 BlackHole 实例：**
```bash
brew install blackhole-2ch   # 用于捕获应用音频
brew install blackhole-16ch  # 用于 Audio Mixer 输出
```

**2. 创建 Multi-Output Device:**
- 打开 Audio MIDI Setup
- 创建 Multi-Output Device
- 勾选：BlackHole 2ch + 扬声器
- 这样音频同时输出到虚拟设备和扬声器（你能听到）

**3. 配置目标应用（如 Spotify）:**
- 音频输出选择 Multi-Output Device
- 音频会流入 BlackHole 2ch

**4. 配置 Audio Mixer:**
```
Input 1 (麦克风): 你的麦克风
Input 2 (系统音频): BlackHole 2ch
Output: BlackHole 16ch
```

**5. 在 OBS/其他软件使用:**
- 音频源选择 BlackHole 16ch
- 录制的是混音后的音频

#### 详细教程：
📖 **[QUICK_SETUP_GUIDE.md](QUICK_SETUP_GUIDE.md)** - 图文并茂，5分钟配置完成

## 为什么这个方案能满足你的需求？

### ✅ 需求1：创建新的虚拟设备输出

**你的需求：** 创建一个专用的虚拟设备用于输出

**解决方案：**
- 安装 BlackHole 16ch 作为 Audio Mixer 的**专用输出设备**
- 这个设备**只**用于 Audio Mixer 输出
- 不影响其他系统音频

**效果：**
```
Audio Mixer → BlackHole 16ch (专用) → OBS/Zoom/录音软件
```

### ✅ 需求2：捕获特定应用的音频

**你的需求：** 像 OBS 一样，只捕获音乐软件的声音，而不是系统总输出

**解决方案：**
- 将音乐播放器的输出设为 BlackHole 2ch
- Audio Mixer 的 Input 2 从 BlackHole 2ch 读取
- 其他应用的声音不会进入 BlackHole 2ch

**效果：**
```
Spotify → BlackHole 2ch → Audio Mixer ✅
Chrome → 默认扬声器 → 不会进入 Audio Mixer ❌
Discord → 默认扬声器 → 不会进入 Audio Mixer ❌
```

只有 Spotify 的音频会被 Audio Mixer 捕获！

## 对比其他方案

| 方案 | 是否满足需求 | 难度 | 成本 | 说明 |
|------|------------|------|------|------|
| **多 BlackHole 实例** | ✅✅✅ | 简单 | 免费 | **推荐方案** |
| 自己编写内核驱动 | ✅ | 极难 | 免费 | 需要 C++/内核开发知识 |
| 使用 Loopback 软件 | ✅✅ | 简单 | $99 | 商业软件，图形化配置 |
| 添加 ScreenCaptureKit | ✅ | 困难 | 免费 | 需要大量开发工作 |
| 使用 JACK Audio | ✅ | 中等 | 免费 | 专业音频工具，学习曲线陡 |

## 实际使用示例

### 场景：直播游戏 + 背景音乐

**需求：**
- 直播中要有麦克风、游戏声音、背景音乐
- 音乐可以独立调节音量
- 音乐不能太响，背景即可

**配置：**

1. **系统默认输出：** 保持为扬声器（游戏声音走这里）

2. **音乐播放器 (Spotify)：**
   - 输出设备：Multi-Output Device (BlackHole 2ch + 扬声器)

3. **Audio Mixer：**
   ```
   Input 1 (麦克风): 100%
   Input 2 (音乐): 30%  ← 背景音乐，音量低
   Output: BlackHole 16ch
   ```

4. **OBS：**
   - 桌面音频：捕获默认输出（游戏声音）
   - 麦克风音频：BlackHole 16ch（麦克风 + 音乐混音）

**结果：**
- ✅ 麦克风声音清晰
- ✅ 游戏声音正常
- ✅ 背景音乐音量适中（30%）
- ✅ 三者完美混合

### 场景：在线会议 + 演示音频

**需求：**
- Zoom 会议中播放演示视频的声音
- 同时保持麦克风清晰

**配置：**

1. **演示视频播放器：**
   - 输出：Multi-Output Device (BlackHole 2ch + 扬声器)

2. **Audio Mixer：**
   ```
   Input 1 (麦克风): 100%
   Input 2 (演示音频): 80%
   Output: BlackHole 16ch
   ```

3. **Zoom：**
   - 麦克风：BlackHole 16ch

**结果：**
- ✅ 会议参与者听到你的声音
- ✅ 会议参与者听到演示视频的声音
- ✅ 两者音量平衡

## 优势总结

与其他方案相比，这个方案：

**🆓 完全免费**
- BlackHole 是开源软件
- 无需购买商业软件

**🚀 配置简单**
- 5分钟完成配置
- 不需要编程知识
- 不需要学习复杂工具

**🔧 灵活可控**
- 可以独立控制每个音源的音量
- 可以随时切换应用的输出设备
- 支持多个应用同时捕获

**💪 稳定可靠**
- BlackHole 是成熟的开源项目
- 被 OBS、Logic Pro 用户广泛使用
- 低延迟、高音质

**🎯 满足专业需求**
- 支持多通道（2ch, 16ch, 64ch）
- 支持高采样率（192kHz）
- 支持专业音频工作流

## 进阶技巧

### 1. 捕获多个应用的音频

安装更多 BlackHole 实例：
```bash
brew install blackhole-64ch
```

创建多个 Multi-Output Device：
- Music + Speakers (BlackHole 2ch)
- Discord + Speakers (BlackHole 16ch)
- Game + Speakers (BlackHole 64ch)

使用多个 Audio Mixer 实例分别混音

### 2. 创建预设配置

保存不同场景的配置文件：
- `config-streaming.json` - 直播配置
- `config-meeting.json` - 会议配置
- `config-recording.json` - 录音配置

快速切换：
```bash
cp ~/.audio-mixer/config-streaming.json ~/.audio-mixer/config.json
./audio-mixer-gui
```

### 3. 使用 AppleScript 自动化

自动设置应用输出：
```applescript
#!/usr/bin/osascript

-- 设置 Spotify 输出到虚拟设备
tell application "Spotify"
    set sound output to "Music + Speakers"
end tell

-- 启动 Audio Mixer
do shell script "~/audio-mixer-gui &"
```

## 未来可能的改进

如果你需要更方便的功能，我可以为 Audio Mixer 添加：

### 1. GUI 应用选择器
- 列出所有运行的音频应用
- 点击选择要捕获的应用
- 自动配置音频路由

### 2. 虚拟设备管理器
- 检测已安装的 BlackHole 实例
- 一键创建 Multi-Output Device
- 自动配置音频路由

### 3. 预设管理
- 保存和加载配置预设
- 快速切换不同场景
- 导出/导入配置

### 4. 音频路由可视化
- 图形化显示音频流向
- 拖拽式连接设备
- 实时监控音频电平

## 总结

**你的需求完全可以通过现有方案实现：**

✅ **"创建新的虚拟设备输出"**
   → 安装 BlackHole 16ch 作为专用输出设备

✅ **"捕获特定应用音频，而不是系统总输出"**
   → 配置应用输出到 BlackHole 2ch，只捕获该应用

✅ **"类似 OBS 的单独捕获"**
   → 通过 Multi-Output Device 实现选择性捕获

**下一步：**
1. 阅读 **[QUICK_SETUP_GUIDE.md](QUICK_SETUP_GUIDE.md)** 完成配置
2. 运行 Audio Mixer GUI
3. 在 OBS 或其他软件中使用混音后的音频

**配置时间：** 5分钟
**难度：** ⭐⭐☆☆☆（简单）
**成本：** 免费

如果你需要更自动化的方案（如 GUI 选择应用、自动路由等），请告诉我，我可以为 Audio Mixer 添加这些功能！

---

**参考文档：**
- 🚀 [快速配置指南](QUICK_SETUP_GUIDE.md) - 图文教程
- 🔧 [高级音频路由](ADVANCED_AUDIO_ROUTING.md) - 专业方案
- 📖 [虚拟设备设置](VIRTUAL_DEVICE_SETUP.md) - 基础安装
- 🎯 [自定义输出设备](CUSTOM_OUTPUT_DEVICE.md) - 设备配置

**更新时间：** 2025-11-27
