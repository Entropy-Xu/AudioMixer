# Audio Mixer Windows Build Script
# Builds both CLI and GUI versions with WASAPI support

param(
    [switch]$Debug,
    [switch]$Release,
    [switch]$SkipDeps
)

Write-Host "=== Audio Mixer Windows Build Script ===" -ForegroundColor Cyan
Write-Host ""

# Function to check if command exists
function Test-Command {
    param([string]$Command)
    try {
        if (Get-Command $Command -ErrorAction Stop) {
            return $true
        }
    }
    catch {
        return $false
    }
}

# Check for Go
if (-not (Test-Command "go")) {
    Write-Host "❌ Error: Go is not installed or not in PATH" -ForegroundColor Red
    Write-Host "   Download from: https://golang.org/dl/" -ForegroundColor Yellow
    exit 1
}

Write-Host "✓ Go found: $(go version)" -ForegroundColor Green
Write-Host ""

# Install dependencies
if (-not $SkipDeps) {
    Write-Host "📦 Downloading dependencies..." -ForegroundColor Cyan
    go mod download
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to download dependencies" -ForegroundColor Red
        exit 1
    }

    go mod tidy
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ Failed to tidy dependencies" -ForegroundColor Red
        exit 1
    }

    Write-Host "✓ Dependencies installed" -ForegroundColor Green
    Write-Host "  ✓ github.com/gordonklaus/portaudio" -ForegroundColor Gray
    Write-Host "  ✓ fyne.io/fyne/v2" -ForegroundColor Gray
    Write-Host "  ✓ github.com/go-ole/go-ole (WASAPI support)" -ForegroundColor Gray
    Write-Host ""
}

# Build CLI version
Write-Host "🔨 Building CLI version..." -ForegroundColor Cyan
go build -o audio-mixer.exe .
if ($LASTEXITCODE -ne 0) {
    Write-Host "❌ CLI build failed" -ForegroundColor Red
    exit 1
}

$cliSize = (Get-Item audio-mixer.exe).Length / 1MB
Write-Host "✓ CLI built: audio-mixer.exe ($([math]::Round($cliSize, 2)) MB)" -ForegroundColor Green
Write-Host ""

# Build GUI version
Write-Host "🖥️  Building GUI version..." -ForegroundColor Cyan

if ($Debug) {
    # Debug build with console window
    Write-Host "   Building DEBUG version (with console)..." -ForegroundColor Yellow
    go build -o audio-mixer-gui.exe ./cmd/gui
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ GUI build failed" -ForegroundColor Red
        exit 1
    }
    $guiSize = (Get-Item audio-mixer-gui.exe).Length / 1MB
    Write-Host "✓ GUI built: audio-mixer-gui.exe ($([math]::Round($guiSize, 2)) MB)" -ForegroundColor Green
    Write-Host "  ⚠ Debug build - console window will be visible" -ForegroundColor Yellow
}
elseif ($Release) {
    # Release build - optimized, no console
    Write-Host "   Building RELEASE version (optimized, no console)..." -ForegroundColor Yellow
    go build -ldflags="-H windowsgui -s -w" -o audio-mixer-gui.exe ./cmd/gui
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ GUI build failed" -ForegroundColor Red
        exit 1
    }
    $guiSize = (Get-Item audio-mixer-gui.exe).Length / 1MB
    Write-Host "✓ GUI built: audio-mixer-gui.exe ($([math]::Round($guiSize, 2)) MB)" -ForegroundColor Green
    Write-Host "  ✓ Optimized: stripped symbols and debug info" -ForegroundColor Gray
    Write-Host "  ✓ Hidden console: no window will appear" -ForegroundColor Gray
}
else {
    # Normal build - no console but not fully optimized
    Write-Host "   Building NORMAL version (no console)..." -ForegroundColor Yellow
    go build -ldflags="-H windowsgui" -o audio-mixer-gui.exe ./cmd/gui
    if ($LASTEXITCODE -ne 0) {
        Write-Host "❌ GUI build failed" -ForegroundColor Red
        exit 1
    }
    $guiSize = (Get-Item audio-mixer-gui.exe).Length / 1MB
    Write-Host "✓ GUI built: audio-mixer-gui.exe ($([math]::Round($guiSize, 2)) MB)" -ForegroundColor Green
    Write-Host "  ✓ Hidden console: no window will appear" -ForegroundColor Gray
}

Write-Host ""
Write-Host "✓ WASAPI support enabled" -ForegroundColor Green
Write-Host "✓ go-ole COM interface integrated" -ForegroundColor Green
Write-Host ""

# Display build summary
Write-Host "🎉 Build complete!" -ForegroundColor Green
Write-Host ""
Write-Host "Built files:" -ForegroundColor Cyan
Write-Host "  • audio-mixer.exe      - CLI version" -ForegroundColor White
Write-Host "  • audio-mixer-gui.exe  - GUI version (with WASAPI)" -ForegroundColor White
Write-Host ""

# Display usage instructions
Write-Host "Usage:" -ForegroundColor Cyan
Write-Host "  Run CLI:  .\audio-mixer.exe" -ForegroundColor White
Write-Host "  Run GUI:  .\audio-mixer-gui.exe" -ForegroundColor White
Write-Host ""

Write-Host "💡 Windows Features:" -ForegroundColor Cyan
Write-Host "  ✓ WASAPI application audio enumeration" -ForegroundColor Green
Write-Host "  ✓ Real-time application list refresh" -ForegroundColor Green
Write-Host "  ✓ Compatible with VB-Cable for app capture" -ForegroundColor Green
Write-Host ""

Write-Host "📖 Next steps:" -ForegroundColor Cyan
Write-Host "  1. Install VB-Cable: https://vb-audio.com/Cable/" -ForegroundColor White
Write-Host "  2. See WINDOWS_SETUP_GUIDE.md for configuration" -ForegroundColor White
Write-Host "  3. See WASAPI_IMPLEMENTATION_NOTES.md for features" -ForegroundColor White
Write-Host ""

# Display build options help
if (-not $Debug -and -not $Release) {
    Write-Host "💡 Build options:" -ForegroundColor Yellow
    Write-Host "  .\build.ps1 -Debug     - Build with console for debugging" -ForegroundColor Gray
    Write-Host "  .\build.ps1 -Release   - Build optimized release version" -ForegroundColor Gray
    Write-Host "  .\build.ps1 -SkipDeps  - Skip dependency download" -ForegroundColor Gray
    Write-Host ""
}
