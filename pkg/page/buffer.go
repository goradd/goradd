package page

import (
	"bytes"
	"github.com/spekary/goradd/pkg/log"
	"sync"
)

// BufferPoolI describes a buffer pool that can be used to improve memory allocation and garbage collection for the
// frequent memory use.
type BufferPoolI interface {
	GetBuffer() *bytes.Buffer
	PutBuffer(buffer *bytes.Buffer)
}

// BufferPool is the global buffer pool used by the page drawing system. You can use it to get buffers for you own
// writes as well. The default buffer pool uses MaxBufferSize to limit the size of buffers that are put back into
// the pool. If a particular http request required a large buffer to satisfy, this prevents that buffer from hanging around too long.
// You should set MaxBufferSize to a value that is bigger than most http request sizes.
var BufferPool BufferPoolI
var MaxBufferSize = 10000

type pool struct {
	sync.Pool
}

func newPool() pool {
	p := pool{
		sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}

	return p
}

func (p pool) GetBuffer() *bytes.Buffer {
	return p.Get().(*bytes.Buffer)
}

func (p pool) PutBuffer(buffer *bytes.Buffer) {
	if buffer.Cap() < MaxBufferSize {
		buffer.Reset()
		p.Put(buffer)
	} else {
		// otherwise we will not put the buffer back, allowing the garbage collector to reclaim the memory
		// Log when our buffer is bigger than MaxBufferSize. If this is happening a lot the value should be increased.
		log.FrameworkDebug("Buffer size was bigger than MaxBufferSize")
	}
}

// GetBuffer returns a buffer from the pool. It will create a new pool if one is not already allocated. This allows
// you to inject your own replacement BufferPool before the first use of GetBuffer()
func GetBuffer() *bytes.Buffer {
	// TODO: If we are running low on memory, notify a sysop.

	if BufferPool == nil {
		BufferPool = newPool()
	}
	return BufferPool.GetBuffer()
}

// PutBuffer puts a buffer back into the pool. Be very careful that you do not refer to the buffer after putting it back,
// including using a slice of a buffer.
func PutBuffer(buffer *bytes.Buffer) {
	if buffer == nil {
		return
	}

	BufferPool.PutBuffer(buffer)
}
