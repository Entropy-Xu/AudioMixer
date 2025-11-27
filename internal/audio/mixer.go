package audio

import (
	"fmt"
	"math"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gordonklaus/portaudio"
)

const (
	// Audio configuration constants
	DefaultSampleRate  = 48000
	DefaultBufferSize  = 512
	DefaultChannels    = 2
	MaxChannels        = 2
	MinLatencyMs       = 10
	MaxLatencyMs       = 100
)

// MixerConfig holds configuration for the audio mixer
type MixerConfig struct {
	SampleRate     float64
	BufferSize     int
	Channels       int
	Input1Device   *portaudio.DeviceInfo // Microphone
	Input2Device   *portaudio.DeviceInfo // Application audio
	OutputDevice   *portaudio.DeviceInfo
	Input1Gain     float32 // 0.0 to 2.0 (0% to 200%)
	Input2Gain     float32
	MasterGain     float32
}

// DefaultMixerConfig returns a default mixer configuration
func DefaultMixerConfig() *MixerConfig {
	return &MixerConfig{
		SampleRate:  DefaultSampleRate,
		BufferSize:  DefaultBufferSize,
		Channels:    DefaultChannels,
		Input1Gain:  1.0,
		Input2Gain:  1.0,
		MasterGain:  1.0,
	}
}

// Mixer handles real-time audio mixing
type Mixer struct {
	config       *MixerConfig
	input1Stream *portaudio.Stream
	input2Stream *portaudio.Stream
	outputStream *portaudio.Stream

	input1Buffer *AudioBuffer
	input2Buffer *AudioBuffer

	bufferPool   *BufferPool

	// Atomic gains for thread-safe volume control
	input1Gain   atomic.Value // float32
	input2Gain   atomic.Value // float32
	masterGain   atomic.Value // float32

	// Metrics
	latency      atomic.Value // time.Duration
	input1Level  atomic.Value // float32
	input2Level  atomic.Value // float32
	outputLevel  atomic.Value // float32

	running      atomic.Bool
	mu           sync.RWMutex
	stopCh       chan struct{}
	wg           sync.WaitGroup
}

// NewMixer creates a new audio mixer
func NewMixer(config *MixerConfig) (*Mixer, error) {
	if config == nil {
		config = DefaultMixerConfig()
	}

	// Validate configuration
	if config.SampleRate <= 0 {
		config.SampleRate = DefaultSampleRate
	}
	if config.BufferSize <= 0 {
		config.BufferSize = DefaultBufferSize
	}
	if config.Channels <= 0 || config.Channels > MaxChannels {
		config.Channels = DefaultChannels
	}

	mixer := &Mixer{
		config:       config,
		input1Buffer: NewAudioBuffer(config.BufferSize * config.Channels * 10),
		input2Buffer: NewAudioBuffer(config.BufferSize * config.Channels * 10),
		bufferPool:   NewBufferPool(config.BufferSize * config.Channels),
		stopCh:       make(chan struct{}),
	}

	// Initialize atomic values
	mixer.input1Gain.Store(config.Input1Gain)
	mixer.input2Gain.Store(config.Input2Gain)
	mixer.masterGain.Store(config.MasterGain)
	mixer.latency.Store(time.Duration(0))
	mixer.input1Level.Store(float32(0))
	mixer.input2Level.Store(float32(0))
	mixer.outputLevel.Store(float32(0))

	return mixer, nil
}

// Start begins audio processing
func (m *Mixer) Start() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.running.Load() {
		return fmt.Errorf("mixer already running")
	}

	// Open input stream 1 (microphone)
	if m.config.Input1Device != nil {
		input1Params := portaudio.StreamParameters{
			Input: portaudio.StreamDeviceParameters{
				Device:   m.config.Input1Device,
				Channels: m.config.Channels,
				Latency:  m.config.Input1Device.DefaultLowInputLatency,
			},
			SampleRate:      m.config.SampleRate,
			FramesPerBuffer: m.config.BufferSize,
		}

		stream, err := portaudio.OpenStream(input1Params, m.input1Callback)
		if err != nil {
			return fmt.Errorf("failed to open input1 stream: %w", err)
		}
		m.input1Stream = stream

		if err := m.input1Stream.Start(); err != nil {
			m.input1Stream.Close()
			return fmt.Errorf("failed to start input1 stream: %w", err)
		}
	}

	// Open input stream 2 (application audio)
	if m.config.Input2Device != nil {
		input2Params := portaudio.StreamParameters{
			Input: portaudio.StreamDeviceParameters{
				Device:   m.config.Input2Device,
				Channels: m.config.Channels,
				Latency:  m.config.Input2Device.DefaultLowInputLatency,
			},
			SampleRate:      m.config.SampleRate,
			FramesPerBuffer: m.config.BufferSize,
		}

		stream, err := portaudio.OpenStream(input2Params, m.input2Callback)
		if err != nil {
			if m.input1Stream != nil {
				m.input1Stream.Close()
			}
			return fmt.Errorf("failed to open input2 stream: %w", err)
		}
		m.input2Stream = stream

		if err := m.input2Stream.Start(); err != nil {
			if m.input1Stream != nil {
				m.input1Stream.Close()
			}
			m.input2Stream.Close()
			return fmt.Errorf("failed to start input2 stream: %w", err)
		}
	}

	// Open output stream
	if m.config.OutputDevice != nil {
		outputParams := portaudio.StreamParameters{
			Output: portaudio.StreamDeviceParameters{
				Device:   m.config.OutputDevice,
				Channels: m.config.Channels,
				Latency:  m.config.OutputDevice.DefaultLowOutputLatency,
			},
			SampleRate:      m.config.SampleRate,
			FramesPerBuffer: m.config.BufferSize,
		}

		stream, err := portaudio.OpenStream(outputParams, m.outputCallback)
		if err != nil {
			if m.input1Stream != nil {
				m.input1Stream.Close()
			}
			if m.input2Stream != nil {
				m.input2Stream.Close()
			}
			return fmt.Errorf("failed to open output stream: %w", err)
		}
		m.outputStream = stream

		if err := m.outputStream.Start(); err != nil {
			if m.input1Stream != nil {
				m.input1Stream.Close()
			}
			if m.input2Stream != nil {
				m.input2Stream.Close()
			}
			m.outputStream.Close()
			return fmt.Errorf("failed to start output stream: %w", err)
		}
	}

	m.running.Store(true)
	return nil
}

// Stop stops audio processing
func (m *Mixer) Stop() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.running.Load() {
		return nil
	}

	m.running.Store(false)
	close(m.stopCh)

	// Stop and close all streams
	var errs []error

	if m.input1Stream != nil {
		if err := m.input1Stream.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("input1 stream stop error: %w", err))
		}
		if err := m.input1Stream.Close(); err != nil {
			errs = append(errs, fmt.Errorf("input1 stream close error: %w", err))
		}
		m.input1Stream = nil
	}

	if m.input2Stream != nil {
		if err := m.input2Stream.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("input2 stream stop error: %w", err))
		}
		if err := m.input2Stream.Close(); err != nil {
			errs = append(errs, fmt.Errorf("input2 stream close error: %w", err))
		}
		m.input2Stream = nil
	}

	if m.outputStream != nil {
		if err := m.outputStream.Stop(); err != nil {
			errs = append(errs, fmt.Errorf("output stream stop error: %w", err))
		}
		if err := m.outputStream.Close(); err != nil {
			errs = append(errs, fmt.Errorf("output stream close error: %w", err))
		}
		m.outputStream = nil
	}

	m.wg.Wait()

	if len(errs) > 0 {
		return fmt.Errorf("errors during stop: %v", errs)
	}

	return nil
}

// input1Callback handles input from first device (microphone)
func (m *Mixer) input1Callback(in []float32) {
	if !m.running.Load() {
		return
	}

	// Calculate and store audio level
	level := calculateRMS(in)
	m.input1Level.Store(level)

	// Write to buffer
	m.input1Buffer.Write(in)
}

// input2Callback handles input from second device (application audio)
func (m *Mixer) input2Callback(in []float32) {
	if !m.running.Load() {
		return
	}

	// Calculate and store audio level
	level := calculateRMS(in)
	m.input2Level.Store(level)

	// Write to buffer
	m.input2Buffer.Write(in)
}

// outputCallback handles output mixing
func (m *Mixer) outputCallback(out []float32) {
	if !m.running.Load() {
		return
	}

	startTime := time.Now()

	// Get buffers from pool
	input1Buf := m.bufferPool.Get()
	input2Buf := m.bufferPool.Get()
	defer func() {
		m.bufferPool.Put(input1Buf)
		m.bufferPool.Put(input2Buf)
	}()

	// Read from input buffers
	m.input1Buffer.Read(input1Buf[:len(out)])
	m.input2Buffer.Read(input2Buf[:len(out)])

	// Get current gains
	input1Gain := m.input1Gain.Load().(float32)
	input2Gain := m.input2Gain.Load().(float32)
	masterGain := m.masterGain.Load().(float32)

	// Mix audio with soft clipping
	for i := range out {
		mixed := (input1Buf[i]*input1Gain + input2Buf[i]*input2Gain) * masterGain
		out[i] = softClip(mixed)
	}

	// Calculate and store output level
	level := calculateRMS(out)
	m.outputLevel.Store(level)

	// Update latency metric
	latency := time.Since(startTime)
	m.latency.Store(latency)
}

// softClip implements soft clipping to prevent harsh distortion
func softClip(sample float32) float32 {
	if sample > 1.0 {
		return 1.0
	}
	if sample < -1.0 {
		return -1.0
	}
	// Soft knee compression near the limits
	if sample > 0.9 {
		return 0.9 + 0.1*float32(math.Tanh(float64(sample-0.9)*5))
	}
	if sample < -0.9 {
		return -0.9 + 0.1*float32(math.Tanh(float64(sample+0.9)*5))
	}
	return sample
}

// calculateRMS calculates the RMS (root mean square) level of audio samples
func calculateRMS(samples []float32) float32 {
	if len(samples) == 0 {
		return 0
	}

	var sum float64
	for _, sample := range samples {
		sum += float64(sample * sample)
	}

	rms := math.Sqrt(sum / float64(len(samples)))
	return float32(rms)
}

// SetInput1Gain sets the gain for input 1 (0.0 to 2.0)
func (m *Mixer) SetInput1Gain(gain float32) {
	if gain < 0 {
		gain = 0
	}
	if gain > 2.0 {
		gain = 2.0
	}
	m.input1Gain.Store(gain)
}

// SetInput2Gain sets the gain for input 2 (0.0 to 2.0)
func (m *Mixer) SetInput2Gain(gain float32) {
	if gain < 0 {
		gain = 0
	}
	if gain > 2.0 {
		gain = 2.0
	}
	m.input2Gain.Store(gain)
}

// SetMasterGain sets the master output gain (0.0 to 2.0)
func (m *Mixer) SetMasterGain(gain float32) {
	if gain < 0 {
		gain = 0
	}
	if gain > 2.0 {
		gain = 2.0
	}
	m.masterGain.Store(gain)
}

// GetInput1Level returns the current RMS level of input 1
func (m *Mixer) GetInput1Level() float32 {
	return m.input1Level.Load().(float32)
}

// GetInput2Level returns the current RMS level of input 2
func (m *Mixer) GetInput2Level() float32 {
	return m.input2Level.Load().(float32)
}

// GetOutputLevel returns the current RMS level of output
func (m *Mixer) GetOutputLevel() float32 {
	return m.outputLevel.Load().(float32)
}

// GetLatency returns the current processing latency
func (m *Mixer) GetLatency() time.Duration {
	return m.latency.Load().(time.Duration)
}

// IsRunning returns whether the mixer is currently running
func (m *Mixer) IsRunning() bool {
	return m.running.Load()
}
