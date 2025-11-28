# Cross-Compiling for Windows

This document explains how to build Windows executables from macOS or Linux.

## ⚠️ Limitations

The Audio Mixer uses:
- **PortAudio** (requires CGO) for cross-platform audio device access
- **Fyne GUI framework** (requires CGO and OpenGL) for the graphical interface

These dependencies make simple cross-compilation challenging because:
1. CGO requires platform-specific C libraries
2. When cross-compiling, the build system can't easily link Windows-specific libraries from macOS/Linux

## ✅ Recommended Solutions

### Option 1: Use fyne-cross (Docker-based)

The easiest way to cross-compile with all dependencies is using `fyne-cross`, which uses Docker containers with pre-configured build environments.

**Requirements:**
- Docker installed and running
- Go installed

**Steps:**

```bash
# Install fyne-cross
go install fyne.io/fyne/v2/cmd/fyne-cross@latest

# Run the cross-compile script (it will detect fyne-cross and use it)
./build-cross-windows.sh
```

This will produce a fully functional Windows executable in `build/audio-mixer-gui-windows-amd64.exe`.

### Option 2: Build Natively on Windows

The most reliable method is to build directly on a Windows machine.

**Requirements:**
- Windows 10/11
- Go installed (https://go.dev/dl/)
- MinGW-w64 or TDM-GCC for CGO support

**Steps:**

1. Clone the repository on Windows
2. Open PowerShell in the project directory
3. Run the build script:

```powershell
.\build.ps1
```

Or manually:

```powershell
go build -o audio-mixer.exe .
go build -ldflags="-H windowsgui -w -s" -o audio-mixer-gui.exe ./cmd/gui
```

See [BUILD_WINDOWS.md](BUILD_WINDOWS.md) for detailed Windows build instructions.

### Option 3: Use GitHub Actions / CI

If you fork this project, you can set up GitHub Actions to automatically build Windows binaries:

```yaml
# .github/workflows/build.yml
name: Build
on: [push, pull_request]
jobs:
  build-windows:
    runs-on: windows-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go build -o audio-mixer.exe .
      - run: go build -ldflags="-H windowsgui" -o audio-mixer-gui.exe ./cmd/gui
      - uses: actions/upload-artifact@v3
        with:
          name: windows-binaries
          path: |
            audio-mixer.exe
            audio-mixer-gui.exe
```

## 🔧 What Doesn't Work

### ❌ Simple mingw-w64 Cross-Compilation

Running this **will not work**:

```bash
# This fails because of CGO dependencies
GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build
```

**Why it fails:**
- PortAudio links to macOS frameworks when building on macOS
- Fyne requires OpenGL headers that differ between platforms
- The linker receives macOS-specific flags like `-framework CoreAudio`

### Partial Solution: Build Without CGO

You can build a limited version without PortAudio:

```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o audio-mixer-limited.exe .
```

**Limitations:**
- No PortAudio device enumeration
- GUI will fail (Fyne requires CGO for OpenGL)
- Only WASAPI code will work on Windows

This is **not recommended** for production use.

## 🎯 Summary

| Method | Difficulty | Result Quality | Requirements |
|--------|-----------|----------------|--------------|
| **fyne-cross** | Easy | ✅ Full features | Docker, Go |
| **Native Windows** | Easy | ✅ Full features | Windows PC |
| **GitHub Actions** | Medium | ✅ Full features | GitHub repo |
| **mingw-w64** | Hard | ❌ Doesn't work | Not viable |

## 📚 Related Documentation

- [BUILD_WINDOWS.md](BUILD_WINDOWS.md) - Native Windows build guide
- [WINDOWS_SETUP_GUIDE.md](WINDOWS_SETUP_GUIDE.md) - Windows usage guide
- [Fyne Cross-Compilation](https://developer.fyne.io/started/cross-compiling)
- [Go CGO Documentation](https://golang.org/cmd/cgo/)

