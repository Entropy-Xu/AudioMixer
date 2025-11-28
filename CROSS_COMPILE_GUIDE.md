# 跨平台编译指南 - 从 Mac/Linux 编译 Windows 版本

本指南说明如何在 macOS 或 Linux 上编译 Windows 的 `.exe` 文件。

---

## 🎯 为什么需要跨平台编译？

- ✅ 在 Mac 上开发，但需要发布 Windows 版本
- ✅ CI/CD 环境中自动构建多平台版本
- ✅ 无需 Windows 虚拟机或双系统
- ✅ 统一的构建流程

---

## 📋 前置要求

### macOS

```bash
# 安装 Homebrew（如果还没有）
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装 mingw-w64 工具链
brew install mingw-w64

# 验证安装
x86_64-w64-mingw32-gcc --version
```

### Ubuntu/Debian

```bash
# 安装 mingw-w64 工具链
sudo apt-get update
sudo apt-get install gcc-mingw-w64 g++-mingw-w64

# 验证安装
x86_64-w64-mingw32-gcc --version
```

### Fedora/RHEL

```bash
# 安装 mingw-w64 工具链
sudo dnf install mingw64-gcc mingw64-gcc-c++

# 验证安装
x86_64-w64-mingw32-gcc --version
```

### Arch Linux

```bash
# 安装 mingw-w64 工具链
sudo pacman -S mingw-w64-gcc

# 验证安装
x86_64-w64-mingw32-gcc --version
```

---

## 🚀 快速开始

### 方法 1: 使用脚本（推荐）

```bash
# 克隆项目
git clone https://github.com/entropy/audio-mixer.git
cd audio-mixer

# 运行交叉编译脚本
./build-cross-windows.sh
```

**输出**:
- `build/audio-mixer-windows-amd64.exe` - CLI 版本
- `build/audio-mixer-gui-windows-amd64.exe` - GUI 版本（含 WASAPI）

### 方法 2: 使用 Makefile

```bash
# 交叉编译 Windows 版本
make cross-windows

# 查看构建结果
ls -lh build/*.exe
```

---

## 📦 构建产物说明

### CLI 版本
- **文件名**: `audio-mixer-windows-amd64.exe`
- **大小**: ~20-30 MB
- **功能**: 命令行音频混音
- **特性**:
  - 完整的音频混音功能
  - PortAudio 支持
  - 命令行交互

### GUI 版本
- **文件名**: `audio-mixer-gui-windows-amd64.exe`
- **大小**: ~30-40 MB
- **功能**: 图形界面音频混音 + WASAPI
- **特性**:
  - ✅ Fyne GUI 界面
  - ✅ WASAPI 应用音频枚举
  - ✅ go-ole COM 接口支持
  - ✅ 隐藏控制台窗口（`-H windowsgui`）
  - ✅ 优化的二进制文件（`-s -w`）

---

## 🔧 技术细节

### 环境变量设置

交叉编译需要设置以下环境变量：

```bash
export GOOS=windows        # 目标操作系统
export GOARCH=amd64        # 目标架构
export CGO_ENABLED=1       # 启用 CGo（必需）
export CC=x86_64-w64-mingw32-gcc    # C 编译器
export CXX=x86_64-w64-mingw32-g++   # C++ 编译器
```

### 编译命令

**CLI 版本**:
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  CXX=x86_64-w64-mingw32-g++ \
  go build -ldflags="-s -w" -o audio-mixer.exe .
```

**GUI 版本（隐藏控制台）**:
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  CXX=x86_64-w64-mingw32-g++ \
  go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui
```

### 编译标志说明

| 标志 | 作用 |
|-----|------|
| `-H windowsgui` | Windows GUI 应用，隐藏控制台窗口 |
| `-s` | 去除符号表（减小文件大小） |
| `-w` | 去除 DWARF 调试信息（减小文件大小） |

---

## ✅ 验证构建

### 检查文件类型

```bash
# 使用 file 命令（macOS/Linux）
file build/audio-mixer-gui-windows-amd64.exe
# 输出应该包含: "PE32+ executable (GUI) x86-64, for MS Windows"
```

### 在 Wine 中测试（可选）

如果安装了 Wine，可以在 macOS/Linux 上运行测试：

```bash
# macOS
brew install wine-stable

# Ubuntu/Debian
sudo apt-get install wine64

# 运行测试
wine build/audio-mixer-gui-windows-amd64.exe
```

⚠️ **注意**: Wine 可能无法完全模拟 Windows 音频环境。

### 在 Windows 上测试（推荐）

将编译的 `.exe` 文件传输到 Windows 机器上测试：

1. 通过网络传输或 USB 驱动器
2. 在 Windows 上双击运行 `audio-mixer-gui-windows-amd64.exe`
3. 测试 WASAPI 功能：
   - 启动音频应用（Spotify、Chrome 等）
   - 点击"🔄 刷新"按钮
   - 验证应用列表是否显示

---

## 🐛 常见问题

### 问题 1: "x86_64-w64-mingw32-gcc: command not found"

**原因**: 未安装 mingw-w64 工具链

**解决方案**: 按照上面的"前置要求"部分安装 mingw-w64

### 问题 2: "undefined reference to PortAudio functions"

**原因**: PortAudio 链接问题

**解决方案**:
```bash
# 清理并重新下载依赖
go clean -cache
go mod download
go mod tidy

# 重新编译
./build-cross-windows.sh
```

### 问题 3: "cannot find package github.com/go-ole/go-ole"

**原因**: go-ole 依赖未下载

**解决方案**:
```bash
go get github.com/go-ole/go-ole
go mod download
```

### 问题 4: 编译成功但 Windows 上无法运行

**可能原因**:
1. 缺少 DLL 依赖
2. 杀毒软件拦截
3. Windows 版本不兼容

**解决方案**:
1. 确保 Windows 10+ 版本
2. 临时禁用杀毒软件测试
3. 使用调试版本查看错误信息：
   ```bash
   # 编译不带 -H windowsgui 的版本
   GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
     CC=x86_64-w64-mingw32-gcc \
     go build -o audio-mixer-gui-debug.exe ./cmd/gui
   ```

### 问题 5: "cgo: C compiler cannot create executables"

**原因**: mingw-w64 配置问题

**解决方案（macOS）**:
```bash
# 重新安装 mingw-w64
brew uninstall mingw-w64
brew install mingw-w64

# 确保在 PATH 中
which x86_64-w64-mingw32-gcc
```

---

## 🎓 进阶用法

### 并行构建多个平台

```bash
# 同时构建所有平台的 GUI 版本
make build-all-gui

# 输出到 build/ 目录
ls -lh build/
```

### 自定义构建标志

编辑 `build-cross-windows.sh` 或直接使用 go build：

```bash
# 添加版本信息
VERSION="v1.0.0"
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 \
  CC=x86_64-w64-mingw32-gcc \
  go build -ldflags="-H windowsgui -s -w -X main.Version=$VERSION" \
  -o audio-mixer-gui-$VERSION.exe ./cmd/gui
```

### 集成到 CI/CD

**GitHub Actions 示例**:

```yaml
name: Cross-compile Windows

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
    - uses: actions/checkout@v3

    - name: Set up Go
      uses: actions/setup-go@v4
      with:
        go-version: '1.21'

    - name: Install mingw-w64
      run: sudo apt-get install -y gcc-mingw-w64

    - name: Cross-compile Windows
      run: make cross-windows

    - name: Upload artifacts
      uses: actions/upload-artifact@v3
      with:
        name: windows-executables
        path: build/*.exe
```

---

## 📊 性能对比

| 构建方式 | 时间 | 文件大小 | 难度 |
|---------|------|---------|------|
| 在 Windows 上构建 | ~30s | 25-35 MB | 简单 |
| 从 Mac 交叉编译 | ~45s | 25-35 MB | 中等 |
| 从 Linux 交叉编译 | ~40s | 25-35 MB | 中等 |

⚠️ 首次编译会更慢（需要下载依赖和编译缓存）

---

## 🎯 最佳实践

### 1. 版本控制

保持构建脚本和 Makefile 在版本控制中：

```bash
git add build-cross-windows.sh Makefile
git commit -m "Add Windows cross-compilation support"
```

### 2. 自动化测试

在 CI 中自动构建并上传：

```yaml
- name: Build and release
  run: |
    make cross-windows
    gh release upload v1.0.0 build/*.exe
```

### 3. 文档更新

确保 README 和文档反映交叉编译选项：

- ✅ 更新 BUILD_QUICK_REFERENCE.md
- ✅ 添加到 README.md 的构建部分
- ✅ 在发布说明中提及

### 4. 依赖管理

定期更新 go.mod 和依赖：

```bash
go get -u ./...
go mod tidy
```

---

## 📚 相关文档

- **[BUILD_QUICK_REFERENCE.md](BUILD_QUICK_REFERENCE.md)** - 快速构建参考
- **[BUILD_WINDOWS.md](BUILD_WINDOWS.md)** - Windows 详细编译指南
- **[WASAPI_IMPLEMENTATION_NOTES.md](WASAPI_IMPLEMENTATION_NOTES.md)** - WASAPI 功能说明
- **[README.md](README.md)** - 项目总览

---

## 🙏 致谢

- **mingw-w64 项目** - 提供 Windows 交叉编译工具链
- **Go 团队** - 优秀的交叉编译支持
- **社区贡献者** - 各种跨平台编译经验分享

---

**最后更新**: 2025-11-27
**适用平台**: macOS, Linux → Windows
**Go 版本**: 1.21+
