# Windows 编译指南

## 快速开始

### 前置条件

1. **安装 Go 1.21+**
   - 下载: https://golang.org/dl/
   - 验证安装: `go version`

2. **安装 PortAudio**
   - Windows 上 PortAudio 会通过 CGo 自动链接
   - 如果遇到问题，可以手动安装 MinGW-w64

3. **安装 git** (如果还没有)
   - 下载: https://git-scm.com/download/win

### 编译步骤

#### 1. 克隆代码库

```bash
git clone https://github.com/entropy/audio-mixer.git
cd audio-mixer
```

#### 2. 安装 Go 依赖

```bash
# 安装所有依赖（包括 WASAPI 支持）
go mod download
```

这会自动下载：
- `github.com/gordonklaus/portaudio` - 音频处理
- `fyne.io/fyne/v2` - GUI 框架
- `github.com/go-ole/go-ole` - Windows COM 接口（WASAPI）

#### 3. 编译 GUI 版本

**普通编译（带控制台窗口）：**
```bash
go build -o audio-mixer-gui.exe ./cmd/gui
```

**无控制台窗口版本（推荐）：**
```bash
go build -ldflags="-H windowsgui" -o audio-mixer-gui.exe ./cmd/gui
```

编译参数说明：
- `-ldflags="-H windowsgui"`: 隐藏控制台窗口
- `-o audio-mixer-gui.exe`: 指定输出文件名
- `./cmd/gui`: GUI 程序的入口目录

#### 4. 运行程序

```bash
./audio-mixer-gui.exe
```

或者直接双击 `audio-mixer-gui.exe`

## WASAPI 功能测试

### 测试应用枚举功能

1. **启动一些音频应用**
   - 播放音乐：Spotify、iTunes、VLC
   - 播放视频：Chrome（YouTube）、Edge
   - 游戏或通讯：Steam、Discord

2. **运行 Audio Mixer**
   ```bash
   ./audio-mixer-gui.exe
   ```

3. **查看应用列表**
   - 在 GUI 中找到"捕获特定应用音频"下拉框
   - 点击"🔄 刷新"按钮
   - 下拉框应该显示所有正在播放音频的应用

4. **验证应用信息**
   - 应用名称应该正确显示（如 "🎵 Spotify"）
   - 应该只显示正在播放音频的应用
   - 切换应用的播放状态后刷新，列表应该更新

### 测试音频捕获（需要 VB-Cable）

要实际捕获应用音频，需要配合 VB-Cable 使用：

1. **安装 VB-Cable**
   - 下载: https://vb-audio.com/Cable/
   - 运行 `VBCABLE_Setup_x64.exe`（或 x86 版本）
   - 安装完成后重启电脑

2. **配置应用输出**
   - 在 Windows 音量混合器中设置目标应用
   - 输出设备选择 "CABLE Input (VB-Audio Virtual Cable)"

3. **配置 Audio Mixer**
   - Input 2: 选择 "CABLE Output"
   - Output: 输入 "CABLE-B Input"（如果安装了多个虚拟设备）

4. **详细配置**
   - 参考: [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md)

## 编译选项

### 调试版本

包含更多日志输出，保留控制台窗口：

```bash
go build -gcflags="-N -l" -o audio-mixer-gui-debug.exe ./cmd/gui
```

### 优化版本

更小的文件大小，更好的性能：

```bash
go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui
```

参数说明：
- `-s`: 去除符号表
- `-w`: 去除 DWARF 调试信息
- 文件大小可减少约 30%

### 静态编译（便携版本）

```bash
CGO_ENABLED=1 go build -ldflags="-H windowsgui -s -w -extldflags '-static'" -o audio-mixer-gui.exe ./cmd/gui
```

注意：静态编译可能需要额外配置 MinGW 工具链

## 常见编译问题

### 问题 1: `cgo: C compiler not found`

**原因**: 缺少 C 编译器（PortAudio 需要 CGo）

**解决方案**:
1. 安装 MinGW-w64
   ```bash
   # 使用 Chocolatey
   choco install mingw

   # 或手动下载
   # https://www.mingw-w64.org/downloads/
   ```

2. 添加到 PATH 环境变量
   ```
   C:\Program Files\mingw-w64\x86_64-8.1.0-posix-seh-rt_v6-rev0\mingw64\bin
   ```

3. 验证安装
   ```bash
   gcc --version
   ```

### 问题 2: `undefined reference to PortAudio functions`

**原因**: PortAudio 链接问题

**解决方案**:
```bash
# 清理缓存并重新编译
go clean -cache
go build -v -o audio-mixer-gui.exe ./cmd/gui
```

### 问题 3: `package github.com/go-ole/go-ole: cannot find package`

**原因**: go-ole 依赖未安装

**解决方案**:
```bash
go get github.com/go-ole/go-ole
go mod tidy
```

### 问题 4: 编译成功但运行时崩溃

**可能原因**:
1. 缺少 DLL 文件
2. 音频驱动问题
3. COM 初始化失败

**解决方案**:
1. 使用调试版本查看错误信息
   ```bash
   go build -o audio-mixer-gui-debug.exe ./cmd/gui
   ./audio-mixer-gui-debug.exe
   ```

2. 检查 Windows 事件查看器的应用程序日志

3. 确保系统有音频设备并正常工作

## 交叉编译（从其他平台编译 Windows 版本）

### 从 macOS/Linux 编译 Windows 版本

```bash
# 设置环境变量
export GOOS=windows
export GOARCH=amd64
export CGO_ENABLED=1

# 安装 MinGW 交叉编译工具链
# macOS:
brew install mingw-w64

# Linux (Ubuntu/Debian):
sudo apt-get install gcc-mingw-w64

# 设置 CGo 编译器
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++

# 编译
go build -ldflags="-H windowsgui" -o audio-mixer-gui.exe ./cmd/gui
```

## 发布版本打包

### 创建发布包

```bash
# 1. 编译优化版本
go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui

# 2. 创建发布目录
mkdir audio-mixer-release
cp audio-mixer-gui.exe audio-mixer-release/
cp README.md audio-mixer-release/
cp WINDOWS_SETUP_GUIDE.md audio-mixer-release/
cp WINDOWS_QUICK_REFERENCE.md audio-mixer-release/
cp WASAPI_IMPLEMENTATION_NOTES.md audio-mixer-release/

# 3. 打包
# 使用 7-Zip 或 Windows 内置压缩
7z a audio-mixer-windows-amd64.zip audio-mixer-release/
```

### 发布清单

发布包应该包含：
- ✅ `audio-mixer-gui.exe` - 主程序
- ✅ `README.md` - 项目说明
- ✅ `WINDOWS_SETUP_GUIDE.md` - Windows 配置指南
- ✅ `WINDOWS_QUICK_REFERENCE.md` - 快速参考
- ✅ `WASAPI_IMPLEMENTATION_NOTES.md` - WASAPI 功能说明
- ✅ `LICENSE` - 开源协议（如果有）

## 性能优化建议

### 编译器优化

使用 Go 编译器的优化选项：

```bash
# 开启所有优化
go build -ldflags="-H windowsgui -s -w" -gcflags="-l=4" -o audio-mixer-gui.exe ./cmd/gui
```

### 运行时优化

设置 Go 运行时环境变量：

```bash
# 减少 GC 压力
set GOGC=200

# 限制 CPU 使用
set GOMAXPROCS=2

# 运行程序
./audio-mixer-gui.exe
```

## 开发工具推荐

### IDE/编辑器

- **VS Code** + Go 扩展
- **GoLand** (JetBrains)
- **Sublime Text** + GoSublime

### 调试工具

- **Delve** - Go 调试器
  ```bash
  go install github.com/go-delve/delve/cmd/dlv@latest
  dlv debug ./cmd/gui
  ```

- **Process Monitor** (Sysinternals)
  - 监控文件和注册表访问
  - 调试 COM 接口调用

- **Dependency Walker**
  - 查看 DLL 依赖
  - 诊断链接问题

## 自动化构建

### 使用 Makefile

创建 `Makefile.windows`:

```makefile
.PHONY: all build clean release

all: build

build:
	go build -o audio-mixer-gui.exe ./cmd/gui

build-release:
	go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui

clean:
	del /Q audio-mixer-gui.exe audio-mixer-gui-debug.exe

release: build-release
	mkdir audio-mixer-release
	copy audio-mixer-gui.exe audio-mixer-release\
	copy README.md audio-mixer-release\
	copy WINDOWS_SETUP_GUIDE.md audio-mixer-release\
	7z a audio-mixer-windows-amd64.zip audio-mixer-release\
```

使用：
```bash
make -f Makefile.windows build-release
```

### 使用 PowerShell 脚本

创建 `build.ps1`:

```powershell
# 编译发布版本
Write-Host "Building release version..." -ForegroundColor Green
go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful!" -ForegroundColor Green

    # 显示文件信息
    $fileInfo = Get-Item audio-mixer-gui.exe
    Write-Host "Size: $($fileInfo.Length / 1MB) MB"
} else {
    Write-Host "Build failed!" -ForegroundColor Red
    exit 1
}
```

使用：
```powershell
.\build.ps1
```

## 持续集成 (CI)

### GitHub Actions 示例

创建 `.github/workflows/build-windows.yml`:

```yaml
name: Build Windows

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: windows-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install dependencies
      run: go mod download

    - name: Build
      run: go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui

    - name: Upload artifact
      uses: actions/upload-artifact@v3
      with:
        name: audio-mixer-windows
        path: audio-mixer-gui.exe
```

## 相关文档

- 📖 [README.md](README.md) - 项目总览
- 🪟 [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md) - Windows 配置指南
- 🔧 [WASAPI_IMPLEMENTATION_NOTES.md](WASAPI_IMPLEMENTATION_NOTES.md) - WASAPI 实现说明
- 📝 [WASAPI_FEATURE_STATUS.md](WASAPI_FEATURE_STATUS.md) - WASAPI 功能状态

## 技术支持

遇到问题？

1. 查看 [常见问题](README.md#故障排除)
2. 查看 [WASAPI 故障排除](WASAPI_IMPLEMENTATION_NOTES.md#常见问题)
3. 提交 GitHub Issue

---

**最后更新**: 2025-11-27
**Go 版本**: 1.21+
**目标平台**: Windows 10/11 (x64)
