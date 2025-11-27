# Windows 配置指南 - 捕获特定应用音频

## 目标

在 Windows 上实现：
- ✅ 捕获特定应用（如音乐播放器）的音频
- ✅ 混合麦克风和应用音频
- ✅ 输出到独立的虚拟设备供其他软件使用

## 快速配置（10分钟）

### 方案选择

| 软件 | 价格 | 难度 | 功能 | 推荐 |
|------|------|------|------|------|
| **VB-Cable** | 免费 | 简单 | 1个虚拟设备 | ⭐⭐⭐ |
| **VB-Audio Hi-Fi Cable** | 免费 | 简单 | 高音质版本 | ⭐⭐⭐⭐ |
| **VB-Audio Matrix** | 免费 | 中等 | 多路混音 | ⭐⭐⭐⭐⭐ |
| **Virtual Audio Cable** | $25 | 中等 | 多个虚拟设备 | ⭐⭐⭐⭐ |

### 推荐配置：VB-Cable + Hi-Fi Cable（免费）

---

## 第1步：安装虚拟音频设备

### 方案 A: VB-Cable（免费，推荐）

#### 1.1 下载 VB-Cable

访问：https://vb-audio.com/Cable/

下载两个软件：
- **VBAUDIO_Cable_Driver.zip** (标准版)
- **VBAUDIO_HFCABLE_Driver.zip** (Hi-Fi 高音质版)

#### 1.2 安装 VB-Cable

1. 解压 `VBAUDIO_Cable_Driver.zip`
2. **右键点击** `VBCABLE_Setup_x64.exe` (64位系统)
3. 选择 **"以管理员身份运行"**
4. 点击 **Install Driver**
5. 等待安装完成

#### 1.3 安装 Hi-Fi Cable

1. 解压 `VBAUDIO_HFCABLE_Driver.zip`
2. **右键点击** `VBCABLE_HiFiCableSetup_x64.exe`
3. 选择 **"以管理员身份运行"**
4. 点击 **Install Driver**
5. 等待安装完成

#### 1.4 重启电脑

**重要**: 必须重启电脑才能正常使用虚拟设备！

### 方案 B: Virtual Audio Cable（付费，$25）

访问：https://vac.muzychenko.net/en/

优势：
- 可以创建多个虚拟设备
- 更稳定
- 更低延迟

---

## 第2步：配置虚拟音频设备

### 2.1 检查虚拟设备

1. **右键点击** 任务栏的音量图标
2. 选择 **"声音"** 或 **"打开声音设置"**
3. 点击 **"声音控制面板"**

应该能看到：
- **播放设备**：
  - CABLE Input (VB-Audio Virtual Cable)
  - CABLE-B Input (VB-Audio Hi-Fi Cable)
- **录制设备**：
  - CABLE Output (VB-Audio Virtual Cable)
  - CABLE-B Output (VB-Audio Hi-Fi Cable)

### 2.2 启用虚拟设备

1. 在 **"播放"** 标签页：
   - 右键空白处 → 勾选 **"显示已禁用的设备"**
   - 右键点击 **CABLE Input** → 选择 **"启用"**
   - 右键点击 **CABLE-B Input** → 选择 **"启用"**

2. 在 **"录制"** 标签页：
   - 右键空白处 → 勾选 **"显示已禁用的设备"**
   - 右键点击 **CABLE Output** → 选择 **"启用"**
   - 右键点击 **CABLE-B Output** → 选择 **"启用"**

### 2.3 配置设备监听（重要！）

**问题**: 如果应用输出到虚拟设备，你将听不到声音。

**解决方案**: 启用"侦听"功能，让虚拟设备的音频转发到扬声器。

#### 配置步骤：

1. 打开 **声音控制面板**
2. 切换到 **"录制"** 标签页
3. 双击 **CABLE Output** 打开属性
4. 切换到 **"侦听"** 标签
5. ✅ 勾选 **"侦听此设备"**
6. 在下拉菜单中选择你的**扬声器/耳机**
7. 点击 **"应用"** 和 **"确定"**

**重复以上步骤** 为 **CABLE-B Output** 配置侦听。

---

## 第3步：配置目标应用音频

### 方法 1: 在应用内设置（推荐）

许多应用支持选择音频输出设备：

#### Spotify:
1. Settings → Audio
2. Output Device → **CABLE Input (VB-Audio Virtual Cable)**

#### Chrome/Edge:
1. 右键点击任务栏音量图标
2. 打开 **音量合成器**
3. 找到 Chrome → 点击设备图标
4. 选择 **CABLE Input**

#### VLC Player:
1. Tools → Preferences
2. Audio → Output module → DirectSound
3. Device → **CABLE Input**

#### Discord:
1. User Settings → Voice & Video
2. Output Device → **CABLE Input**

### 方法 2: 使用 Windows 应用音量和设备首选项（Windows 10/11）

1. 右键点击任务栏音量图标
2. 选择 **"打开音量合成器"** 或 **"声音设置"**
3. 向下滚动到 **"应用音量和设备首选项"**
4. 找到目标应用（如 Spotify）
5. 在 **"输出"** 下拉菜单中选择 **CABLE Input**

**优势**:
- 无需在应用内设置
- 统一管理所有应用的音频输出
- 可以为不同应用设置不同的虚拟设备

---

## 第4步：配置 Audio Mixer

### 4.1 编译 Windows 版本

```cmd
# 在项目目录
go build -o audio-mixer-gui.exe ./cmd/gui
```

### 4.2 运行 Audio Mixer

```cmd
audio-mixer-gui.exe
```

### 4.3 配置设备

在 GUI 中设置：

```
Input 1 (麦克风):
  选择你的麦克风设备（如 Realtek High Definition Audio）

Input 2 (系统音频):
  选择 <Auto Detect Loopback>
  或手动选择: [Virtual] CABLE Output

Output (虚拟输出设备名称):
  输入: CABLE-B Input
  点击 "检测设备" 验证
```

### 4.4 调节音量

```
Input 1 (麦克风): 100%
Input 2 (应用音频): 70-80%
Master: 100%
```

### 4.5 开始混音

点击 **"Start Mixer"**

---

## 第5步：在其他应用中使用混音

### OBS Studio

1. **Sources** → **Add** → **Audio Input Capture**
2. 创建新源，命名为 "Mixed Audio"
3. **Device** → 选择 **CABLE-B Output (VB-Audio Hi-Fi Cable)**
4. 点击 **OK**

现在 OBS 会录制麦克风 + 应用音频的混音。

### Zoom / Teams / Skype

1. 打开 **Audio Settings**
2. **Microphone** → 选择 **CABLE-B Output**
3. 参与会议时，你的声音 + 应用音频会被传输

### Audacity / Adobe Audition

1. **Edit** → **Preferences** → **Devices**
2. **Recording Device** → 选择 **CABLE-B Output**
3. 点击录音按钮

---

## 完整音频流图（Windows）

```
┌─────────────────────────────┐
│ Spotify / Chrome / Discord  │
│ 输出设置: CABLE Input       │
└──────────┬──────────────────┘
           │
           ▼
   ┌───────────────────┐
   │ CABLE Input       │ (虚拟播放设备)
   └─────────┬─────────┘
             │
             ▼
   ┌───────────────────┐
   │ CABLE Output      │ (虚拟录制设备)
   └─────────┬─────────┘
             │
             ├────────────────────────┐
             │                        │
             ▼                        ▼
   ┌──────────────────┐     ┌──────────────────┐
   │ Audio Mixer      │     │ 扬声器/耳机      │
   │ Input 2          │     │ (侦听功能)       │
   └─────────┬────────┘     └──────────────────┘
             │                        ▲
   ┌─────────┴────────┐              │
   │ 麦克风            │              │ 你能听到音乐
   │ Input 1          │              │
   └─────────┬────────┘              │
             │                        │
             ▼                        │
   ┌────────────────────────┐        │
   │ 混音处理                │        │
   │ (麦克风 + 应用音频)     │        │
   └─────────┬──────────────┘        │
             │                        │
             ▼                        │
   ┌────────────────────┐            │
   │ CABLE-B Input      │            │
   └─────────┬──────────┘            │
             │                        │
             ▼                        │
   ┌────────────────────┐            │
   │ CABLE-B Output     │────────────┘
   └─────────┬──────────┘
             │
             ▼
   ┌────────────────────────────┐
   │ OBS / Zoom / 录音软件       │
   │ 音频输入: CABLE-B Output    │
   └────────────────────────────┘
```

---

## 实际使用场景

### 场景 1: 直播游戏 + 背景音乐

**配置:**

1. **游戏**:
   - 音频输出: 默认扬声器（不被捕获）

2. **音乐播放器 (Spotify)**:
   - 输出设备: **CABLE Input**

3. **Audio Mixer**:
   ```
   Input 1: 麦克风 (100%)
   Input 2: CABLE Output (30% - 背景音乐)
   Output: CABLE-B Input
   ```

4. **OBS**:
   - 桌面音频: 默认扬声器（游戏声音）
   - 麦克风: CABLE-B Output（麦克风 + 音乐混音）

**结果:**
- ✅ 游戏声音完整
- ✅ 麦克风清晰
- ✅ 背景音乐音量适中

### 场景 2: 播客录制 + 音效

**配置:**

1. **音效播放器**:
   - 输出: **CABLE Input**

2. **Audio Mixer**:
   ```
   Input 1: 麦克风 (100%)
   Input 2: CABLE Output (50% - 音效)
   Output: CABLE-B Input
   ```

3. **Audacity**:
   - 录音设备: CABLE-B Output

**结果:**
- 录制包含麦克风和音效的播客

### 场景 3: 在线会议 + 演示视频

**配置:**

1. **演示视频播放器**:
   - 输出: **CABLE Input**

2. **Audio Mixer**:
   ```
   Input 1: 麦克风 (100%)
   Input 2: CABLE Output (80% - 演示音频)
   Output: CABLE-B Input
   ```

3. **Zoom / Teams**:
   - 麦克风: CABLE-B Output

**结果:**
- 会议参与者听到你的声音和演示视频的声音

---

## 高级配置

### 使用多个虚拟设备

如果你需要分别控制多个应用的音频：

1. **安装 Virtual Audio Cable (VAC)** - 支持创建多个虚拟设备

2. **创建多个虚拟线路**:
   - Line 1: 音乐播放器
   - Line 2: Discord
   - Line 3: 游戏音频

3. **使用多个 Audio Mixer 实例**

### 使用 Voicemeeter（免费替代方案）

**Voicemeeter** 是一个功能强大的虚拟混音台：

下载：https://vb-audio.com/Voicemeeter/

**特点:**
- 图形化混音界面
- 多输入多输出
- 内置均衡器和压缩器
- 完全免费

**与 Audio Mixer 对比:**

| 功能 | Audio Mixer | Voicemeeter |
|------|-------------|-------------|
| 开源 | ✅ | ❌ |
| 跨平台 | ✅ | ❌ (仅 Windows) |
| 简单易用 | ✅ | ⚠️ (学习曲线) |
| 音效处理 | ❌ | ✅ |
| 图形化路由 | ❌ | ✅ |

**建议:**
- 如果只需要简单混音 → 使用 Audio Mixer
- 如果需要复杂音频路由 → 使用 Voicemeeter

---

## 故障排除

### ❌ 问题 1: 听不到应用的声音

**原因:** 应用输出到虚拟设备后，没有配置侦听

**解决方案:**
1. 打开 **声音控制面板**
2. **录制** 标签页 → **CABLE Output** → **属性**
3. **侦听** 标签 → ✅ **侦听此设备**
4. 选择你的扬声器
5. **应用** 和 **确定**

### ❌ 问题 2: OBS 录制没声音

**检查项:**
1. ✅ Audio Mixer 是否正在运行？
2. ✅ OBS 音频源是否选择了 CABLE-B Output？
3. ✅ Input 2 是否有音频波形？
4. ✅ Windows 音量合成器中，应用是否静音？

**解决方案:**
- 在 Audio Mixer 中查看 Input 2 的电平表
- 确保应用正在播放音频
- 检查应用的输出设备是否是 CABLE Input

### ❌ 问题 3: 音频延迟高

**解决方案:**
1. 降低 Audio Mixer 的 **Buffer Size**:
   - 512 → 256 samples
2. 关闭其他音频应用
3. 在 **声音控制面板** 中:
   - CABLE Output → 属性 → 高级
   - 降低独占模式的采样率

### ❌ 问题 4: 音质下降 / 有杂音

**解决方案:**
1. 提高 **Buffer Size**: 256 → 512 或 1024
2. 统一采样率:
   - 所有设备使用相同的采样率（48000 Hz）
3. 在 **声音控制面板** 中:
   - 右键设备 → 属性 → 高级
   - 设为: **24 位，48000 Hz**

### ❌ 问题 5: VB-Cable 驱动无法安装

**原因:** Windows 驱动签名验证

**解决方案:**
1. **以管理员身份运行** 安装程序
2. 如果仍然失败:
   - 重启电脑进入 **安全模式**
   - 禁用驱动签名验证
   - 安装 VB-Cable
   - 重启到正常模式

### ❌ 问题 6: Audio Mixer 提示 "Invalid number of channels"

**解决方案:**
- 这个错误已在最新代码中修复
- 确保使用最新版本的 Audio Mixer
- 如果问题仍存在，尝试将配置文件中的 `channels` 设为 2

---

## 推荐配置参数

### 直播/游戏场景
```json
{
  "sample_rate": 48000,
  "buffer_size": 512,
  "channels": 2,
  "input1_gain": 1.0,
  "input2_gain": 0.3,
  "master_gain": 1.0
}
```

### 播客/录音场景
```json
{
  "sample_rate": 48000,
  "buffer_size": 1024,
  "channels": 2,
  "input1_gain": 1.0,
  "input2_gain": 0.5,
  "master_gain": 1.0
}
```

### 在线会议场景
```json
{
  "sample_rate": 44100,
  "buffer_size": 1024,
  "channels": 2,
  "input1_gain": 1.0,
  "input2_gain": 0.8,
  "master_gain": 1.0
}
```

---

## 总结

**Windows 上实现你的需求的步骤:**

1. ✅ **安装虚拟设备** - VB-Cable + Hi-Fi Cable (免费)
2. ✅ **配置侦听** - 让你能听到虚拟设备的音频
3. ✅ **设置应用输出** - 目标应用输出到 CABLE Input
4. ✅ **配置 Audio Mixer** - Input 2 从 CABLE Output 读取，Output 到 CABLE-B Input
5. ✅ **在 OBS/Zoom 使用** - 音频源选择 CABLE-B Output

**配置时间:** 10分钟
**难度:** ⭐⭐☆☆☆
**成本:** 免费 (使用 VB-Cable)

**下一步:**
- 按照本指南完成配置
- 测试麦克风和应用音频是否都能被捕获
- 在 OBS 或其他软件中验证混音效果

如果遇到问题，参考 **故障排除** 章节！

---

**相关文档:**
- 📖 [虚拟设备设置指南](VIRTUAL_DEVICE_SETUP.md) - 跨平台虚拟设备安装
- 🚀 [快速配置指南](QUICK_SETUP_GUIDE.md) - macOS 版本
- 🔧 [高级音频路由](ADVANCED_AUDIO_ROUTING.md) - 专业音频路由方案

**更新时间:** 2025-11-27
