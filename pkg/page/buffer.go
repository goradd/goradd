package page

import (
	"bytes"
	"github.com/goradd/goradd/pkg/log"
	"sync"
)

// BufferPoolI describes a buffer pool that can be used to improve memory allocation and garbage collection for
// frequently used memory buffers.
type BufferPoolI interface {
	GetBuffer() *bytes.Buffer
	PutBuffer(buffer *bytes.Buffer)
}

// BufferPool is the global buffer pool used by the page drawing system. You can use it to get buffers for you own
// writes as well. The default buffer pool uses MaxBufferSize to limit the size of buffers that are put back into
// the pool. If a particular http request required a large buffer to satisfy, this prevents that buffer from hanging around too long.
// You should set MaxBufferSize to a value that is bigger than most http request sizes.
var BufferPool BufferPoolI

// MaxBufferSize is the maximum size that a buffer will be allowed to grow before it is automatically removed
// from the buffer pool. This prevents large memory allocations from permanently sitting in the buffer pool.
// You should set MaxBufferSize to a value that is bigger than most http request sizes. The default is 10,000 bytes.
var MaxBufferSize = 10000

type pool struct {
	sync.Pool
}

// TODO: Test and improve the allocation mechanism here under heavy load. We could potentially run out of memory
// so we should attempt to limit how much memory the pool is allowed to hold on to. The sync.Pool documentation
// says that the fmt package has an example of how to use pool such that it scales under heavy load, but
// releases memory when quiet.

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

// GetBuffer returns a buffer from the pool if one is available, or creates a new one if all the buffers are being used.
// Generally, you should follow a GetBuffer with a deferred PutBuffer, and the PutBuffer should be in the same
// function as the GetBuffer to prevent memory leaks.
func (p pool) GetBuffer() *bytes.Buffer {
	return p.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the buffer pool. Always do this after you are done with a buffer. If the buffer
// has grown bigger than MaxBufferSize, it will be removed from the pool so that it can be garbage collected.
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
