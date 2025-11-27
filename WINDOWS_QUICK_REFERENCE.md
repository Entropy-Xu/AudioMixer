# Windows 快速参考卡片

## 一页配置清单 ✓

### 准备工作 (一次性)

#### 1. 安装虚拟音频驱动

```
下载: https://vb-audio.com/Cable/
├─ VBAUDIO_Cable_Driver.zip       (标准版)
└─ VBAUDIO_HFCABLE_Driver.zip     (Hi-Fi版)

安装步骤:
1. 解压文件
2. 右键 → 以管理员身份运行
3. Install Driver
4. 重启电脑 ⚠️
```

#### 2. 启用侦听功能

```
声音控制面板 → 录制:
├─ CABLE Output → 属性 → 侦听
│   └─ ✅ 侦听此设备 → 选择扬声器
└─ CABLE-B Output → 属性 → 侦听
    └─ ✅ 侦听此设备 → 选择扬声器
```

### 每次使用配置

#### 3. 设置应用输出

**方法A: 应用内设置**
```
Spotify: Settings → Output Device → CABLE Input
Chrome: 音量合成器 → Chrome → CABLE Input
Discord: Audio Settings → Output → CABLE Input
```

**方法B: Windows 设置 (推荐)**
```
设置 → 系统 → 声音 → 应用音量和设备首选项
└─ 找到应用 → 输出 → CABLE Input
```

#### 4. 配置 Audio Mixer

```
Input 1 (麦克风): 你的麦克风
Input 2 (系统音频): <Auto Detect Loopback>
Output: CABLE-B Input
```

#### 5. 在 OBS/Zoom 使用

```
OBS: 音频源 → CABLE-B Output
Zoom: 麦克风 → CABLE-B Output
```

---

## 音频流向图

```
应用 → CABLE Input → CABLE Output → Audio Mixer Input 2
                          ↓
                      扬声器 (侦听)

麦克风 → Audio Mixer Input 1

Audio Mixer → CABLE-B Input → CABLE-B Output → OBS/Zoom
```

---

## 设备对照表

| 虚拟设备 | 用途 | 配置位置 |
|----------|------|----------|
| **CABLE Input** | 应用音频输出到这里 | 应用内设置 |
| **CABLE Output** | Audio Mixer 从这里读取 | Audio Mixer Input 2 |
| **CABLE-B Input** | Audio Mixer 输出到这里 | Audio Mixer Output |
| **CABLE-B Output** | OBS/Zoom 从这里读取 | OBS/Zoom 音频源 |

---

## 常见问题速查

| 问题 | 解决方案 |
|------|----------|
| 听不到应用声音 | 检查 CABLE Output 的"侦听"是否启用 |
| OBS 没声音 | 确认 Audio Mixer 正在运行 + OBS 选择 CABLE-B Output |
| 音频延迟高 | 降低 Buffer Size (512 → 256) |
| 音质差/杂音 | 提高 Buffer Size (256 → 1024) + 统一采样率 48000Hz |
| 驱动安装失败 | 以管理员身份运行 + 重启电脑 |

---

## 推荐配置参数

### 直播/游戏
```
Sample Rate: 48000 Hz
Buffer Size: 512 samples
Input 1 (麦克风): 100%
Input 2 (音乐): 30%
Master: 100%
```

### 播客/录音
```
Sample Rate: 48000 Hz
Buffer Size: 1024 samples
Input 1 (麦克风): 100%
Input 2 (音效): 50%
Master: 100%
```

### 在线会议
```
Sample Rate: 44100 Hz
Buffer Size: 1024 samples
Input 1 (麦克风): 100%
Input 2 (演示音频): 80%
Master: 100%
```

---

## 使用场景示例

### 🎮 直播游戏 + 背景音乐

```
1. Spotify → 输出: CABLE Input
2. 游戏 → 输出: 默认扬声器
3. Audio Mixer:
   Input 1: 麦克风 (100%)
   Input 2: CABLE Output (30%)
   Output: CABLE-B Input
4. OBS:
   桌面音频: 默认 (游戏声)
   麦克风: CABLE-B Output (麦克风+音乐)
```

### 🎙️ 播客录制 + 音效

```
1. 音效播放器 → 输出: CABLE Input
2. Audio Mixer:
   Input 1: 麦克风 (100%)
   Input 2: CABLE Output (50%)
   Output: CABLE-B Input
3. Audacity:
   录音设备: CABLE-B Output
```

### 💼 在线会议 + 演示视频

```
1. 视频播放器 → 输出: CABLE Input
2. Audio Mixer:
   Input 1: 麦克风 (100%)
   Input 2: CABLE Output (80%)
   Output: CABLE-B Input
3. Zoom/Teams:
   麦克风: CABLE-B Output
```

---

## 故障排除命令

### 检查虚拟设备状态

打开 PowerShell 执行:
```powershell
# 列出所有音频设备
Get-CimInstance Win32_SoundDevice | Select-Object Name, Status
```

### 重置虚拟设备

```
1. 设备管理器
2. 声音、视频和游戏控制器
3. 找到 VB-Audio Virtual Cable
4. 右键 → 禁用设备 → 启用设备
```

### 完全卸载重装

```
1. 运行安装程序
2. 选择 Uninstall
3. 重启电脑
4. 重新安装
5. 再次重启
```

---

## 软件下载链接

| 软件 | 链接 | 价格 |
|------|------|------|
| VB-Cable | https://vb-audio.com/Cable/ | 免费 |
| Virtual Audio Cable | https://vac.muzychenko.net/en/ | $25 |
| Voicemeeter | https://vb-audio.com/Voicemeeter/ | 免费 |
| Audio Mixer (本项目) | https://github.com/your-repo | 开源免费 |

---

## 键盘快捷键 (Audio Mixer)

| 操作 | 快捷键 |
|------|--------|
| 启动混音器 | 点击 Start Mixer |
| 停止混音器 | 点击 Stop Mixer |
| 检测设备 | 点击 检测设备 |

*注: 可以为 Audio Mixer 添加全局快捷键支持*

---

## 性能优化建议

### 降低 CPU 使用率
```
1. 提高 Buffer Size (512 → 1024)
2. 关闭不必要的音频应用
3. 禁用音效增强 (声音控制面板 → 增强)
```

### 降低延迟
```
1. 降低 Buffer Size (1024 → 256)
2. 使用独占模式 (设备属性 → 高级)
3. 提高进程优先级 (任务管理器 → Audio Mixer → 高)
```

### 提升音质
```
1. 统一采样率: 48000 Hz
2. 统一位深度: 24 bit
3. 禁用所有音效增强
4. 使用高品质虚拟驱动 (Hi-Fi Cable)
```

---

## 更新日志

**v2.1 - 2025-11-27**
- ✅ 添加自定义输出设备名称功能
- ✅ 修复通道数错误
- ✅ 完善 Windows 配置文档

**v2.0 - 2025-11-27**
- ✅ 虚拟设备支持
- ✅ 系统音频捕获
- ✅ GUI 界面优化

---

**需要帮助？**
- 📖 完整文档: [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md)
- 💬 问题反馈: GitHub Issues
- 🚀 快速开始: 跟随本页清单配置

**预计配置时间:** 10-15 分钟
**难度等级:** ⭐⭐☆☆☆ (简单)
