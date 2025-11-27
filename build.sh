#!/bin/bash

set -e

echo "=== Audio Mixer Build Script ==="
echo ""

# Get dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod tidy

echo ""
echo "🔨 Building CLI version..."
go build -o audio-mixer .
echo "✅ CLI version built: ./audio-mixer"

echo ""
echo "🖥️  Building GUI version..."
# Suppress duplicate library warnings on macOS
if [[ "$OSTYPE" == "darwin"* ]]; then
    go build -ldflags="-w -s" -o audio-mixer-gui ./cmd/gui 2>&1 | grep -v "ignoring duplicate libraries" || true
else
    go build -ldflags="-w -s" -o audio-mixer-gui ./cmd/gui
fi
echo "✅ GUI version built: ./audio-mixer-gui"

echo ""
echo "🎉 Build complete!"
echo ""
echo "Run CLI: ./audio-mixer"
echo "Run GUI: ./audio-mixer-gui"
