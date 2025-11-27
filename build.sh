#!/bin/bash

set -e

echo "=== Audio Mixer Build Script ==="
echo ""

# Detect platform
PLATFORM=$(uname -s)
echo "🖥️  Platform detected: $PLATFORM"
echo ""

# Get dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "🔨 Building CLI version..."
if [[ "$PLATFORM" == "MINGW"* ]] || [[ "$PLATFORM" == "MSYS"* ]] || [[ "$PLATFORM" == "CYGWIN"* ]]; then
    # Windows
    go build -o audio-mixer.exe .
    echo "✅ CLI version built: ./audio-mixer.exe"
else
    # macOS/Linux
    go build -o audio-mixer .
    echo "✅ CLI version built: ./audio-mixer"
fi

echo ""
echo "🖥️  Building GUI version..."
# Suppress duplicate library warnings on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    go build -ldflags="-w -s" -o audio-mixer-gui ./cmd/gui 2>&1 | grep -v "ignoring duplicate libraries" || true
    echo "✅ GUI version built: ./audio-mixer-gui"
elif [[ "$PLATFORM" == "MINGW"* ]] || [[ "$PLATFORM" == "MSYS"* ]] || [[ "$PLATFORM" == "CYGWIN"* ]]; then
    # Windows with hidden console
    echo "   Building with WASAPI support..."
    go build -ldflags="-H windowsgui -w -s" -o audio-mixer-gui.exe ./cmd/gui
    echo "✅ GUI version built: ./audio-mixer-gui.exe"
    echo "   ✓ WASAPI application capture enabled"
else
    # Linux
    go build -ldflags="-w -s" -o audio-mixer-gui ./cmd/gui
    echo "✅ GUI version built: ./audio-mixer-gui"
fi

echo ""
echo "🎉 Build complete!"
echo ""
if [[ "$PLATFORM" == "MINGW"* ]] || [[ "$PLATFORM" == "MSYS"* ]] || [[ "$PLATFORM" == "CYGWIN"* ]]; then
    echo "Run CLI: ./audio-mixer.exe"
    echo "Run GUI: ./audio-mixer-gui.exe"
    echo ""
    echo "💡 Windows Features:"
    echo "   - WASAPI application audio enumeration"
    echo "   - Use with VB-Cable for application capture"
    echo "   - See WINDOWS_SETUP_GUIDE.md for details"
else
    echo "Run CLI: ./audio-mixer"
    echo "Run GUI: ./audio-mixer-gui"
fi
