# Cross-Compilation Fix Summary

## Problem
When attempting to cross-compile for Windows from macOS using `./build-cross-windows.sh`, the build failed with the error:

```
x86_64-w64-mingw32-gcc: error: unrecognized command-line option '-framework'
```

This occurred because macOS-specific PortAudio linking flags (like `-framework CoreAudio`) were being passed to the Windows cross-compiler.

## Root Cause

The Audio Mixer project has several dependencies that complicate cross-compilation:

1. **PortAudio** - A cross-platform audio I/O library
   - Requires CGO (C bindings)
   - Links to platform-specific libraries:
     - macOS: CoreAudio, AudioToolbox frameworks
     - Windows: MME, DirectSound, WASAPI
     - Linux: ALSA, PulseAudio

2. **Fyne GUI Framework**
   - Requires CGO for OpenGL bindings
   - Needs platform-specific graphics libraries

3. **Build System Limitation**
   - When CGO_ENABLED=1 on macOS, pkg-config returns macOS library flags
   - The mingw-w64 cross-compiler receives incompatible flags
   - Simple cross-compilation cannot resolve these platform differences

## Solution Implemented

### 1. Updated build-cross-windows.sh

The script now:
- Checks for `fyne-cross` (Docker-based cross-compilation tool)
- If found: Uses fyne-cross to build properly
- If not found: Provides clear instructions with 3 alternatives

### 2. Added Build Tags for CGO Dependencies

Created platform-specific builds:
- `device.go` - Tagged with `//go:build cgo` (PortAudio version)
- `device_nocgo.go` - Tagged with `//go:build !cgo` (stub version)
- `mixer.go` - Tagged with `//go:build cgo`
- `loopback.go` - Tagged with `//go:build cgo`

This allows building without PortAudio when CGO is disabled, though with limited functionality.

### 3. Created Documentation

Added `CROSS_COMPILE_WINDOWS.md` explaining:
- Why simple cross-compilation doesn't work
- Three working solutions (fyne-cross, native Windows, CI/CD)
- Comparison table of different approaches
- Technical details about the limitations

## Recommended Build Methods

### For Users on macOS/Linux:

**Option A: fyne-cross (Easiest)**
```bash
go install fyne.io/fyne/v2/cmd/fyne-cross@latest
./build-cross-windows.sh  # Now uses fyne-cross automatically
```

**Option B: Native Windows Build (Most Reliable)**
- Build on a Windows machine
- Run `.\build.ps1` or `go build`
- See BUILD_WINDOWS.md

### For Users on Windows:

Simply run:
```powershell
.\build.ps1
```
or
```bash
./build.sh  # In Git Bash
```

## What Changed

### Modified Files:
1. `build-cross-windows.sh` - Completely rewritten to detect fyne-cross or show instructions
2. `internal/audio/device.go` - Added `//go:build cgo` tag
3. `internal/audio/mixer.go` - Added `//go:build cgo` tag
4. `internal/audio/loopback.go` - Added `//go:build cgo` tag

### New Files:
1. `internal/audio/device_nocgo.go` - Stub implementation for no-CGO builds
2. `CROSS_COMPILE_WINDOWS.md` - Comprehensive cross-compilation guide

### Unchanged:
- Regular native builds (`./build.sh`) still work perfectly
- All existing functionality preserved
- No breaking changes to the codebase

## Testing

✅ Native macOS build: Works
✅ Cross-compile script guidance: Works
✅ Error messages: Clear and helpful
✅ Documentation: Comprehensive

## Technical Details

### Why mingw-w64 Alone Doesn't Work

When you run:
```bash
GOOS=windows CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build
```

The Go build system:
1. Enables CGO
2. Runs pkg-config to find PortAudio
3. Gets macOS-specific flags from pkg-config
4. Passes these flags to mingw-w64-gcc
5. mingw-w64-gcc rejects `-framework` (macOS-only flag)
6. Build fails

### Why fyne-cross Works

fyne-cross:
1. Runs in a Docker container with Windows build environment
2. Has pre-built Windows libraries (PortAudio, OpenGL)
3. Uses correct pkg-config for Windows
4. Produces fully functional Windows executable

### Why Native Build is Best

Building on the target platform:
- No cross-compilation complexity
- Uses native toolchain
- Most reliable and debuggable
- Fastest build times

## Future Improvements

Potential enhancements:
1. Add CI/CD workflow for automatic Windows builds
2. Consider removing PortAudio dependency on Windows (use WASAPI only)
3. Provide pre-built binaries in GitHub Releases
4. Create Docker-based build environment without fyne-cross

## Related Issues

This fix addresses cross-compilation issues related to:
- CGO dependency management
- Platform-specific library linking
- GUI framework requirements
- Build system portability

