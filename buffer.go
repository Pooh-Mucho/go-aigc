package aigc

import (
	"bytes"
	"sync"
)

var bufferPool = sync.Pool{New: func() interface{} { return &bytes.Buffer{} }}

// AllocBuffer returns a buffer from the bufferPool.
func AllocBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

// FreeBuffer returns a buffer to the bufferPool.
func FreeBuffer(buf *bytes.Buffer) {
	if buf == nil {
		return
	}
	if buf.Len() > 64*1024 {
		return
	}
	buf.Reset()
	bufferPool.Put(buf)
}
