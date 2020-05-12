package gst

import (
	"bytes"
)

// BytePool implements a leaky pool of []byte in the form of a bounded
// channel. It has to be for the similar arrays to avoid memory wasting
type BytePool struct {
	c chan []byte
	w int
}

// NewBytePool creates a new BytePool bounded to the given maxSize, with new
// byte arrays sized based on width.
func NewBytePool(maxSize int, width int) (bp *BytePool) {
	return &BytePool{
		c: make(chan []byte, maxSize),
		w: width,
	}
}

// Get gets a []byte from the BytePool, or creates a new one if none are
// available in the pool.
func (bp *BytePool) Get() (b []byte) {
	select {
	case b = <-bp.c:
	// reuse existing buffer
	default:
		// create new buffer
		b = make([]byte, bp.w)
		//log.Println("Make new one!")
	}
	return
}

// Put returns the given Buffer to the BytePool.
func (bp *BytePool) Put(b []byte) {
	if cap(b) < bp.w {
		// someone tried to put back a too small buffer, discard it
		return
	} else if cap(b) > bp.w {
		//since we expect to work with raw frames mostly
		//the size of them should be the same
		bp.w = cap(b)
	}

	select {
	case bp.c <- b:
		// buffer went back into pool
	default:
		// buffer didn't go back into pool, just discard
	}
}

// NumPooled returns the number of items currently pooled.
func (bp *BytePool) NumPooled() int {
	return len(bp.c)
}

// Width returns the width of the byte arrays in this pool.
func (bp *BytePool) Width() (n int) {
	return bp.w
}

// BufferPool implements a pool of bytes.Buffers in the form of a bounded
// channel.
type BufferPool struct {
	c chan *bytes.Buffer
}

// NewBufferPool creates a new BufferPool bounded to the given size.
func NewBufferPool(size int) (bp *BufferPool) {
	return &BufferPool{
		c: make(chan *bytes.Buffer, size),
	}
}

// Get gets a Buffer from the BufferPool, or creates a new one if none are
// available in the pool.
func (bp *BufferPool) Get() (b *bytes.Buffer) {
	select {
	case b = <-bp.c:
	// reuse existing buffer
	default:
		// create new buffer
		b = bytes.NewBuffer([]byte{})
	}
	return
}

// Put returns the given Buffer to the BufferPool.
func (bp *BufferPool) Put(b *bytes.Buffer) {
	b.Reset()
	select {
	case bp.c <- b:
	default: // Discard the buffer if the pool is full.
	}
}

// NumPooled returns the number of items currently pooled.
func (bp *BufferPool) NumPooled() int {
	return len(bp.c)
}
