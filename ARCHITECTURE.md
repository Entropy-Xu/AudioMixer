# 架构设计文档

## 概述

Audio Mixer是一个实时低延迟音频混音工具,采用Go语言开发,使用PortAudio作为跨平台音频接口。

## 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Main Process                         │
│  ┌─────────────┐  ┌──────────────┐  ┌──────────────────┐   │
│  │   Config    │  │    Device    │  │      Mixer       │   │
│  │   Manager   │  │    Manager   │  │     Engine       │   │
│  └─────────────┘  └──────────────┘  └──────────────────┘   │
└─────────────────────────────────────────────────────────────┘
                            │
        ┌───────────────────┼───────────────────┐
        │                   │                   │
        ▼                   ▼                   ▼
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│   Input 1    │   │   Input 2    │   │    Output    │
│  (Callback)  │   │  (Callback)  │   │  (Callback)  │
└──────────────┘   └──────────────┘   └──────────────┘
        │                   │                   ▲
        │                   │                   │
        ▼                   ▼                   │
┌──────────────┐   ┌──────────────┐            │
│ Ring Buffer 1│   │ Ring Buffer 2│            │
└──────────────┘   └──────────────┘            │
        │                   │                   │
        └───────────┬───────────────────────────┘
                    │
                    ▼
            ┌──────────────┐
            │  Mix + Gain  │
            │  Soft Clip   │
            └──────────────┘
```

## 核心组件

### 1. Device Manager (`internal/audio/device.go`)

**职责**: 音频设备管理

**功能**:
- 初始化PortAudio
- 枚举所有音频设备
- 筛选输入/输出设备
- 获取设备详细信息
- 管理设备生命周期

**接口**:
```go
type DeviceManager struct {
    initialized bool
}

func (dm *DeviceManager) Initialize() error
func (dm *DeviceManager) Terminate() error
func (dm *DeviceManager) ListDevices() ([]*DeviceInfo, error)
func (dm *DeviceManager) GetInputDevices() ([]*DeviceInfo, error)
func (dm *DeviceManager) GetOutputDevices() ([]*DeviceInfo, error)
```

### 2. Audio Buffer (`internal/audio/buffer.go`)

**职责**: 内存管理和音频数据缓冲

**组件**:

#### BufferPool
- 使用 `sync.Pool` 减少GC压力
- 复用音频缓冲区
- 零分配设计

#### AudioBuffer (Ring Buffer)
- 无锁循环缓冲区
- 线程安全的读写
- 固定大小,防止内存增长

**性能优化**:
```go
// 从池中获取buffer,使用后归还
buf := bufferPool.Get()
defer bufferPool.Put(buf)
```

### 3. Mixer Engine (`internal/audio/mixer.go`)

**职责**: 核心混音处理

**配置**:
```go
type MixerConfig struct {
    SampleRate     float64  // 采样率(默认48000Hz)
    BufferSize     int      // 缓冲区大小(默认512帧)
    Channels       int      // 声道数(默认2)
    Input1Device   *DeviceInfo
    Input2Device   *DeviceInfo
    OutputDevice   *DeviceInfo
    Input1Gain     float32  // 0.0 - 2.0
    Input2Gain     float32
    MasterGain     float32
}
```

**处理流程**:

1. **Input Callbacks** (音频线程):
   ```go
   func input1Callback(in []float32) {
       // 1. 计算RMS电平
       level := calculateRMS(in)

       // 2. 存储电平(atomic)
       m.input1Level.Store(level)

       // 3. 写入ring buffer
       m.input1Buffer.Write(in)
   }
   ```

2. **Output Callback** (音频线程):
   ```go
   func outputCallback(out []float32) {
       // 1. 从ring buffer读取
       m.input1Buffer.Read(input1Buf)
       m.input2Buffer.Read(input2Buf)

       // 2. 获取增益(atomic)
       gain1 := m.input1Gain.Load()
       gain2 := m.input2Gain.Load()

       // 3. 混音
       for i := range out {
           mixed := (input1[i]*gain1 + input2[i]*gain2) * masterGain
           out[i] = softClip(mixed)
       }

       // 4. 计算输出电平
       m.outputLevel.Store(calculateRMS(out))
   }
   ```

### 4. Config Manager (`internal/config/config.go`)

**职责**: 配置管理和持久化

**配置文件位置**: `~/.audio-mixer/config.json`

**功能**:
- 加载/保存配置
- 配置验证
- 默认值管理

## 数据流

### 音频数据流
```
Hardware Mic → PortAudio → Input1Callback → Ring Buffer 1 ─┐
                                                             │
Hardware App → PortAudio → Input2Callback → Ring Buffer 2 ─┤
                                                             │
                                                             ├→ Mix → Output
                                                             │
Volume Controls (atomic) ────────────────────────────────────┘
```

### 控制流
```
User Input → Config → Mixer Config → Atomic Values → Audio Callbacks
```

## 线程模型

### 线程类型

1. **主线程**:
   - UI/CLI交互
   - 配置管理
   - 设备管理
   - 启动/停止混音器

2. **PortAudio线程** (实时优先级):
   - Input1 callback
   - Input2 callback
   - Output callback
   - 不能阻塞!

3. **监控线程**:
   - 定期读取电平和延迟
   - 更新UI显示

### 线程安全

- **Atomic操作**: 音量控制和指标
- **Ring Buffer**: 内部使用RWMutex
- **Buffer Pool**: sync.Pool是线程安全的
- **无共享可变状态**: 各callback独立

## 混音算法

### 简单加法混音
```go
output = (input1 * gain1 + input2 * gain2) * masterGain
```

### 软削波(Soft Clipping)
```go
func softClip(sample float32) float32 {
    // 硬限制
    if sample > 1.0 { return 1.0 }
    if sample < -1.0 { return -1.0 }

    // 软膝压缩(0.9以上)
    if sample > 0.9 {
        // 使用tanh实现平滑过渡
        return 0.9 + 0.1*tanh((sample-0.9)*5)
    }
    if sample < -0.9 {
        return -0.9 + 0.1*tanh((sample+0.9)*5)
    }

    return sample
}
```

**原理**:
- 防止削波失真
- 在接近限制时平滑压缩
- 比硬限制听感更好

### RMS电平计算
```go
func calculateRMS(samples []float32) float32 {
    var sum float64
    for _, sample := range samples {
        sum += float64(sample * sample)
    }
    return sqrt(sum / len(samples))
}
```

## 性能优化

### 1. 内存管理

**问题**: 频繁分配/释放导致GC暂停

**解决方案**:
```go
// 使用sync.Pool
bufferPool := NewBufferPool(size)
buf := bufferPool.Get()    // 从池获取
defer bufferPool.Put(buf)  // 归还池
```

### 2. 无锁通信

**问题**: 锁竞争导致延迟

**解决方案**:
```go
// 使用atomic操作
m.input1Gain.Store(newGain)     // 设置
gain := m.input1Gain.Load()     // 读取
```

### 3. Ring Buffer

**问题**: 队列操作开销大

**解决方案**:
- 固定大小循环缓冲区
- 预分配内存
- 简单的指针操作

### 4. 零拷贝

**问题**: 数据拷贝浪费CPU

**解决方案**:
```go
// 直接操作callback提供的buffer
func outputCallback(out []float32) {
    // 不拷贝,直接写入out
    for i := range out {
        out[i] = mixed[i]
    }
}
```

## 延迟分析

### 延迟组成

```
Total Latency = Input Latency + Buffer Latency + Processing Latency + Output Latency
```

1. **Input Latency**: 硬件 + 驱动 (~5-10ms)
2. **Buffer Latency**: `bufferSize / sampleRate`
   - 512帧 @ 48kHz = 10.7ms
   - 1024帧 @ 48kHz = 21.3ms
3. **Processing Latency**: 混音计算 (~0.1-1ms)
4. **Output Latency**: 硬件 + 驱动 (~5-10ms)

**总延迟**: 通常 20-30ms

### 降低延迟

1. 减小buffer size (但可能导致音频断续)
2. 使用ASIO/CoreAudio低延迟驱动
3. 优化混音算法
4. 提高进程优先级

## 错误处理

### 策略

1. **初始化错误**: 返回error,由调用者处理
2. **运行时错误**: 记录日志,尝试恢复
3. **配置错误**: 使用默认值,警告用户

### 错误类型

```go
// 致命错误 - 停止程序
if err := portaudio.Initialize(); err != nil {
    return fmt.Errorf("fatal: %w", err)
}

// 非致命错误 - 降级处理
if input2Device == nil {
    log.Warn("Input 2 not configured, running mono-mix mode")
}
```

## 扩展点

### 1. 音频效果器

在混音前/后插入效果器:

```go
// Pre-mix effects
input1Processed := applyEQ(input1, eqSettings)

// Mix
mixed := mix(input1Processed, input2)

// Post-mix effects
output := applyCompressor(mixed, compSettings)
```

### 2. 多路输入

扩展为N路输入:

```go
type Mixer struct {
    inputs []InputChannel
}

type InputChannel struct {
    stream *Stream
    buffer *RingBuffer
    gain   atomic.Value
}
```

### 3. 录音功能

在output callback中写入文件:

```go
func outputCallback(out []float32) {
    // Mix audio
    mixed := mix(...)

    // Write to file
    if recording {
        recordBuffer.Write(mixed)
    }

    copy(out, mixed)
}
```

## 平台特定实现

### macOS
- CoreAudio backend
- 支持聚合设备(Aggregate Device)
- 低延迟模式

### Windows
- WASAPI backend
- ASIO支持(需要驱动)
- DirectSound fallback

### Linux
- PulseAudio backend
- ALSA直接访问
- JACK支持(专业音频)

## 测试策略

### 单元测试
```bash
go test ./internal/audio -v
go test ./internal/config -v
```

### 集成测试
- Loopback测试(输出连回输入)
- 延迟测量
- 长时间稳定性测试

### 性能测试
```bash
go test -bench=. ./internal/audio
```

## 监控和调试

### 指标
- Input/Output电平(RMS)
- 处理延迟
- CPU使用率
- 内存分配

### 调试技巧
```go
// 启用PortAudio调试
// export PA_DEBUG=1

// Go race detector
go run -race main.go

// CPU profiling
go test -cpuprofile=cpu.prof
go tool pprof cpu.prof
```

## 未来改进

1. **SIMD优化**: 使用SIMD指令加速混音
2. **GPU加速**: 音频FFT处理
3. **机器学习**: 智能降噪和音质增强
4. **插件系统**: 支持VST/AU效果器
5. **网络同步**: 多机器协同混音

## 参考资源

- [PortAudio文档](http://www.portaudio.com/docs/v19-doxydocs/)
- [Go音频处理最佳实践](https://github.com/golang/go/wiki/AudioProgramming)
- [实时音频编程指南](https://www.rossbencina.com/code/real-time-audio-programming-101-time-waits-for-nothing)
