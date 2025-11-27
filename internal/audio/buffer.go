package audio

import (
	"sync"
)

// BufferPool manages audio buffer allocation to reduce GC pressure
type BufferPool struct {
	pool sync.Pool
	size int
}

// NewBufferPool creates a new buffer pool with specified buffer size
func NewBufferPool(size int) *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				buf := make([]float32, size)
				return &buf
			},
		},
		size: size,
	}
}

// Get retrieves a buffer from the pool
func (bp *BufferPool) Get() []float32 {
	bufPtr := bp.pool.Get().(*[]float32)
	return *bufPtr
}

// Put returns a buffer to the pool
func (bp *BufferPool) Put(buf []float32) {
	// Clear the buffer before returning to pool
	for i := range buf {
		buf[i] = 0
	}
	bp.pool.Put(&buf)
}

// AudioBuffer represents a thread-safe ring buffer for audio samples
type AudioBuffer struct {
	data     []float32
	readPos  int
	writePos int
	size     int
	mu       sync.RWMutex
}

// NewAudioBuffer creates a new ring buffer
func NewAudioBuffer(size int) *AudioBuffer {
	return &AudioBuffer{
		data: make([]float32, size),
		size: size,
	}
}

// Write writes samples to the buffer
func (ab *AudioBuffer) Write(samples []float32) int {
	ab.mu.Lock()
	defer ab.mu.Unlock()

	written := 0
	for _, sample := range samples {
		ab.data[ab.writePos] = sample
		ab.writePos = (ab.writePos + 1) % ab.size
		written++
	}
	return written
}

// Read reads samples from the buffer
func (ab *AudioBuffer) Read(samples []float32) int {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	read := 0
	for i := range samples {
		samples[i] = ab.data[ab.readPos]
		ab.readPos = (ab.readPos + 1) % ab.size
		read++
	}
	return read
}

// Available returns the number of samples available to read
func (ab *AudioBuffer) Available() int {
	ab.mu.RLock()
	defer ab.mu.RUnlock()

	if ab.writePos >= ab.readPos {
		return ab.writePos - ab.readPos
	}
	return ab.size - ab.readPos + ab.writePos
}
