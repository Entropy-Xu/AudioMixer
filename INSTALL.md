# 安装指南

## 系统依赖

### macOS

1. 安装Homebrew(如果尚未安装):
```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

2. 安装PortAudio:
```bash
brew install portaudio
```

3. 安装Go(如果尚未安装):
```bash
brew install go
```

4. 验证安装:
```bash
go version
portaudio --version || echo "PortAudio installed"
```

### Linux (Ubuntu/Debian)

1. 更新包管理器:
```bash
sudo apt-get update
```

2. 安装依赖:
```bash
sudo apt-get install -y portaudio19-dev build-essential
```

3. 安装Go:
```bash
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

4. 验证安装:
```bash
go version
```

### Linux (Fedora/RHEL)

```bash
sudo dnf install portaudio-devel gcc
```

### Windows

1. 安装Go: 从 [golang.org](https://golang.org/dl/) 下载安装器

2. 安装MinGW-w64(用于CGO):
   - 下载 [MinGW-w64](https://www.mingw-w64.org/)
   - 或使用 [MSYS2](https://www.msys2.org/):
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```

3. 添加MinGW到PATH环境变量

4. PortAudio会通过CGO自动处理

## 构建项目

1. 克隆或下载项目:
```bash
cd audio-mixer
```

2. 下载Go依赖:
```bash
go mod download
```

3. 构建:
```bash
go build -o audio-mixer .
```

或使用Makefile:
```bash
make build
```

4. 运行:
```bash
./audio-mixer
```

## 虚拟音频设备设置

### macOS - BlackHole

1. 安装BlackHole:
```bash
brew install blackhole-2ch
```

或从 [GitHub](https://github.com/ExistentialAudio/BlackHole/releases) 下载安装包

2. 创建聚合设备(Aggregate Device):
   - 打开 `/Applications/Utilities/Audio MIDI Setup.app`
   - 点击左下角 `+` 按钮
   - 选择 "Create Aggregate Device"
   - 勾选你的实际输出设备和BlackHole
   - 命名为 "Audio Mixer Output"

3. 创建多输出设备(Multi-Output Device):
   - 再次点击 `+`
   - 选择 "Create Multi-Output Device"
   - 勾选你想同时输出的设备
   - 用于将应用音频同时发送到BlackHole和扬声器

### Windows - VB-Cable

1. 下载VB-Cable:
   - 访问 [VB-Audio官网](https://vb-audio.com/Cable/)
   - 下载免费版本

2. 解压并右键点击 `VBCABLE_Setup_x64.exe`
   - 选择 "以管理员身份运行"
   - 按提示完成安装

3. 重启计算机

4. 在Windows声音设置中:
   - 将应用程序输出设置为 "CABLE Input"
   - Audio Mixer的Input 2选择 "CABLE Output"

### Linux - PulseAudio

1. 列出所有音频源(包括monitor):
```bash
pacmd list-sources | grep -e 'name:' -e 'device.description'
```

2. 创建虚拟sink(可选):
```bash
pactl load-module module-null-sink sink_name=VirtualSink sink_properties=device.description="Virtual_Sink"
```

3. 使用pavucontrol GUI管理:
```bash
sudo apt-get install pavucontrol
pavucontrol
```

## 权限设置

### macOS

首次运行时,系统会请求麦克风访问权限:
- 系统偏好设置 → 安全性与隐私 → 隐私 → 麦克风
- 确保 Terminal(或你的终端应用)有麦克风权限

### Linux

确保你的用户在audio组:
```bash
sudo usermod -a -G audio $USER
```

重新登录使更改生效。

## 故障排除

### 编译错误: "portaudio.h: No such file or directory"

**原因**: PortAudio库未安装或未在系统PATH中

**解决方案**:
- macOS: `brew install portaudio`
- Linux: `sudo apt-get install portaudio19-dev`
- Windows: 确保MinGW已安装

### 运行时错误: "Failed to initialize PortAudio"

**原因**: 音频系统未正确配置

**解决方案**:
- 检查音频设备是否被其他程序占用
- 重启音频服务
- macOS: `sudo killall coreaudiod`
- Linux: `pulseaudio --kill && pulseaudio --start`

### macOS: "operation not permitted"

**原因**: 麦克风权限未授予

**解决方案**:
1. 系统偏好设置 → 安全性与隐私 → 隐私 → 麦克风
2. 添加并勾选你的终端应用
3. 重启终端和程序

### 延迟过高

**原因**: 缓冲区太大

**解决方案**:
编辑 `~/.audio-mixer/config.json`:
```json
{
  "buffer_size": 256
}
```

较小的值 = 更低延迟,但可能导致音频断断续续
建议范围: 256-1024

### 音频断续/爆音

**原因**: 缓冲区太小或CPU过载

**解决方案**:
1. 增大buffer_size:
```json
{
  "buffer_size": 1024
}
```

2. 关闭其他占用CPU的应用
3. 检查系统资源使用情况

## 验证安装

运行以下命令验证一切正常:

```bash
./audio-mixer
```

应该看到:
- 列出所有可用音频设备
- 可以选择输入/输出设备
- 启动后显示实时音频电平

如果看到设备列表,说明安装成功!

## 下一步

阅读 [README.md](README.md) 了解:
- 使用指南
- 配置选项
- 使用场景示例
