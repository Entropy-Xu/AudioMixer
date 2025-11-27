// +build ignore

package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/entropy/audio-mixer/internal/audio"
)

// 简单的音频直通示例
// 将麦克风输入直接传递到输出设备
func main() {
	fmt.Println("=== Simple Audio Passthrough Example ===")
	fmt.Println("This example captures microphone input and plays it back to the output device")
	fmt.Println()

	// 初始化设备管理器
	deviceManager := audio.NewDeviceManager()
	if err := deviceManager.Initialize(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer deviceManager.Terminate()

	// 获取默认设备
	inputDevice, err := deviceManager.GetDefaultInputDevice()
	if err != nil {
		fmt.Printf("Error getting input device: %v\n", err)
		os.Exit(1)
	}

	outputDevice, err := deviceManager.GetDefaultOutputDevice()
	if err != nil {
		fmt.Printf("Error getting output device: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Input Device:  %s\n", inputDevice.Name)
	fmt.Printf("Output Device: %s\n", outputDevice.Name)
	fmt.Println()

	// 创建混音器配置(只使用一路输入)
	config := audio.DefaultMixerConfig()
	config.Input1Device = inputDevice
	config.Input2Device = nil // 不使用第二路输入
	config.OutputDevice = outputDevice
	config.Input1Gain = 1.0 // 100%音量
	config.MasterGain = 1.0

	// 创建混音器
	mixer, err := audio.NewMixer(config)
	if err != nil {
		fmt.Printf("Error creating mixer: %v\n", err)
		os.Exit(1)
	}

	// 启动混音器
	fmt.Println("Starting audio passthrough...")
	if err := mixer.Start(); err != nil {
		fmt.Printf("Error starting mixer: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ Passthrough active!")
	fmt.Println("Speak into your microphone - you should hear yourself")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// 监控音频电平
	ticker := time.NewTicker(200 * time.Millisecond)
	defer ticker.Stop()

	stopCh := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				level := mixer.GetInput1Level()
				db := levelToDB(level)
				bar := getLevelBar(level, 40)
				fmt.Printf("\rInput Level: %6.1f dB [%s]", db, bar)
			case <-stopCh:
				return
			}
		}
	}()

	// 等待中断信号
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	<-sigCh

	fmt.Println("\n\nStopping...")
	close(stopCh)

	// 停止混音器
	if err := mixer.Stop(); err != nil {
		fmt.Printf("Error stopping mixer: %v\n", err)
	}

	fmt.Println("Goodbye!")
}

func levelToDB(level float32) float32 {
	if level < 0.00001 {
		return -100.0
	}
	return 20.0 * float32(math.Log10(float64(level)))
}

func getLevelBar(level float32, width int) string {
	filled := int(level * float32(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}
	return strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
}
