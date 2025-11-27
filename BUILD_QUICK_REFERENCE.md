# 构建快速参考

快速查找如何在各平台上构建 Audio Mixer。

---

## 🪟 Windows

### 使用 PowerShell 脚本（推荐）

```powershell
# 普通构建（无控制台窗口）
.\build.ps1

# 调试构建（带控制台窗口，方便查看错误）
.\build.ps1 -Debug

# 发布构建（完全优化，文件最小）
.\build.ps1 -Release

# 跳过依赖下载（加快构建）
.\build.ps1 -SkipDeps
```

**输出**:
- `audio-mixer.exe` - CLI 版本
- `audio-mixer-gui.exe` - GUI 版本（含 WASAPI 支持）

### 使用 Git Bash / MSYS2

```bash
./build.sh
```

自动检测 Windows 环境并构建 `.exe` 文件。

### 使用 Makefile

```bash
# 仅构建 Windows 版本
make build-windows

# 构建 Windows 调试版本
make build-windows-debug

# 查看所有选项
make help
```

**输出位置**: `build/` 目录

---

## 🍎 macOS

### 使用 Shell 脚本（推荐）

```bash
./build.sh
```

**输出**:
- `audio-mixer` - CLI 版本
- `audio-mixer-gui` - GUI 版本

### 使用 Makefile

```bash
# 构建 CLI
make build

# 构建 GUI
make gui

# 构建两者（优化版本）
make build-release

# 查看所有选项
make help
```

---

## 🐧 Linux

### 使用 Shell 脚本

```bash
./build.sh
```

### 使用 Makefile

```bash
# 构建 CLI
make build

# 构建 GUI
make gui

# 查看所有选项
make help
```

---

## 🌍 跨平台编译

### 从任何平台构建所有平台版本

**构建所有平台的 CLI 版本**:
```bash
make build-all
```

**构建所有平台的 GUI 版本**:
```bash
make build-all-gui
```

**输出**: `build/` 目录
- `audio-mixer-darwin-amd64` - macOS Intel
- `audio-mixer-darwin-arm64` - macOS Apple Silicon
- `audio-mixer-linux-amd64` - Linux x64
- `audio-mixer-windows-amd64.exe` - Windows x64 (含 WASAPI)

---

## 📋 常用命令速查

| 需求 | Windows (PowerShell) | macOS/Linux | Makefile |
|-----|---------------------|-------------|----------|
| **快速构建** | `.\build.ps1` | `./build.sh` | `make gui` |
| **调试版本** | `.\build.ps1 -Debug` | `./build.sh` | `make build` |
| **发布版本** | `.\build.ps1 -Release` | `./build.sh` | `make build-release` |
| **仅 Windows** | `.\build.ps1` | - | `make build-windows` |
| **所有平台** | - | - | `make build-all-gui` |
| **清理** | `Remove-Item *.exe` | `make clean` | `make clean` |

---

## 🔧 构建选项说明

### Windows 构建标志

| 标志 | 用途 | 何时使用 |
|-----|------|---------|
| `-H windowsgui` | 隐藏控制台窗口 | 发布版本 |
| `-s` | 去除符号表 | 发布版本（减小文件） |
| `-w` | 去除调试信息 | 发布版本（减小文件） |
| 无标志 | 保留调试信息 | 开发和调试 |

### 文件大小对比

| 版本 | 文件大小（约） |
|-----|--------------|
| Debug（完整信息） | ~40-50 MB |
| Normal（无控制台） | ~30-40 MB |
| Release（完全优化） | ~20-30 MB |

---

## ⚡ 快速开始

### Windows 新用户

1. 安装 Go: https://golang.org/dl/
2. 克隆仓库: `git clone https://github.com/entropy/audio-mixer.git`
3. 进入目录: `cd audio-mixer`
4. 构建: `.\build.ps1`
5. 运行: `.\audio-mixer-gui.exe`

### macOS 新用户

1. 安装依赖:
   ```bash
   brew install go portaudio
   ```
2. 克隆仓库: `git clone https://github.com/entropy/audio-mixer.git`
3. 进入目录: `cd audio-mixer`
4. 构建: `./build.sh`
5. 运行: `./audio-mixer-gui`

### Linux 新用户

1. 安装依赖:
   ```bash
   sudo apt-get install golang portaudio19-dev
   ```
2. 克隆仓库: `git clone https://github.com/entropy/audio-mixer.git`
3. 进入目录: `cd audio-mixer`
4. 构建: `./build.sh`
5. 运行: `./audio-mixer-gui`

---

## 🐛 故障排除

### Windows: "cgo: C compiler not found"

**解决方案**: 安装 MinGW-w64
```powershell
# 使用 Chocolatey
choco install mingw

# 或手动下载
# https://www.mingw-w64.org/downloads/
```

### macOS: "ld: library not found for -lportaudio"

**解决方案**: 安装 PortAudio
```bash
brew install portaudio
```

### Linux: "portaudio.h: No such file or directory"

**解决方案**: 安装开发包
```bash
sudo apt-get install portaudio19-dev
```

### 所有平台: "package github.com/go-ole/go-ole: cannot find package"

**解决方案**: 下载依赖
```bash
go mod download
go mod tidy
```

---

## 📦 发布打包

### 创建 Windows 发布包

```powershell
# 1. 构建发布版本
.\build.ps1 -Release

# 2. 创建发布目录
New-Item -ItemType Directory -Force -Path release
Copy-Item audio-mixer-gui.exe release/
Copy-Item README.md release/
Copy-Item WINDOWS_SETUP_GUIDE.md release/
Copy-Item WASAPI_IMPLEMENTATION_NOTES.md release/

# 3. 打包
Compress-Archive -Path release\* -DestinationPath audio-mixer-windows-v1.0.0.zip
```

### 创建 macOS 发布包

```bash
# 1. 构建发布版本
./build.sh

# 2. 创建 .app 包（可选）
# 或直接打包二进制文件
mkdir -p release
cp audio-mixer-gui release/
cp README.md release/
cp QUICK_SETUP_GUIDE.md release/

# 3. 打包
tar -czf audio-mixer-macos-v1.0.0.tar.gz release/
```

---

## 🎯 特定场景

### 场景 1: 开发 Windows 功能

```powershell
# 使用调试版本，可以看到控制台输出
.\build.ps1 -Debug
.\audio-mixer-gui.exe
```

### 场景 2: 准备发布

```bash
# 构建所有平台的优化版本
make build-all-gui

# 检查输出
ls -lh build/
```

### 场景 3: 快速迭代开发

```bash
# 使用 Makefile 的 watch 模式（如果有）
# 或手动快速构建
make gui && ./audio-mixer-gui
```

### 场景 4: CI/CD 集成

```yaml
# GitHub Actions 示例
- name: Build
  run: |
    go mod download
    make build-all-gui
```

---

## 📚 相关文档

- **详细构建指南**: [BUILD_WINDOWS.md](BUILD_WINDOWS.md)
- **Windows 配置**: [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md)
- **WASAPI 功能**: [WASAPI_IMPLEMENTATION_NOTES.md](WASAPI_IMPLEMENTATION_NOTES.md)
- **项目说明**: [README.md](README.md)

---

**最后更新**: 2025-11-27
