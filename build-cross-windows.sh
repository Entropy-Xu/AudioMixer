#!/bin/bash

# Cross-compile Windows executables from macOS/Linux
# Requires mingw-w64 toolchain

set -e

echo "=== Cross-compile Windows from macOS/Linux ==="
echo ""

# Check for fyne-cross
if command -v fyne-cross &> /dev/null; then
    echo "✓ fyne-cross found - using Docker-based cross-compilation"
    echo ""

    # Get dependencies
    echo "📦 Downloading dependencies..."
    go mod download
    go mod tidy

    mkdir -p build

    echo ""
    echo "🖥️  Cross-compiling GUI for Windows using fyne-cross..."
    fyne-cross windows -arch=amd64 -app-id=com.github.entropy.audiomixer ./cmd/gui

    # Move the output to our build directory
    if [ -f "fyne-cross/bin/windows-amd64/gui.exe" ]; then
        mv fyne-cross/bin/windows-amd64/gui.exe build/audio-mixer-gui-windows-amd64.exe
        echo "✓ GUI built: build/audio-mixer-gui-windows-amd64.exe"
    fi

    echo ""
    echo "🎉 Cross-compilation complete!"
    echo ""
    ls -lh build/*.exe 2>/dev/null
    exit 0
fi

echo "❌ Cross-compilation limitation detected"
echo ""
echo "The Audio Mixer GUI uses Fyne which requires CGO and OpenGL."
echo "Simple cross-compilation is not possible without additional tools."
echo ""
echo "📖 See CROSS_COMPILE_WINDOWS.md for detailed explanation"
echo ""
echo "Options to build for Windows:"
echo ""
echo "1. Use fyne-cross (Docker-based cross-compilation):"
echo "   go install fyne.io/fyne/v2/cmd/fyne-cross@latest"
echo "   Then run this script again"
echo ""
echo "2. Build natively on Windows:"
echo "   • Install Go on Windows"
echo "   • Run: .\\build.ps1"
echo "   • See BUILD_WINDOWS.md for details"
echo ""
echo "3. Use GitHub Actions / CI:"
echo "   • The project can build Windows binaries in CI"
echo "   • Check the releases page for pre-built binaries"
echo ""
exit 1
