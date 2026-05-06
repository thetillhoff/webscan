package webserver

import (
	"bytes"
	"sync"
)

type OutputBuffer struct {
	buffer bytes.Buffer
	mu     sync.Mutex
}

func NewOutputBuffer() *OutputBuffer {
	return &OutputBuffer{}
}

func (ob *OutputBuffer) Write(p []byte) (n int, err error) {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.buffer.Write(p)
}

func (ob *OutputBuffer) Bytes() []byte {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.buffer.Bytes()
}

func (ob *OutputBuffer) String() string {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	return ob.buffer.String()
}

func (ob *OutputBuffer) Reset() {
	ob.mu.Lock()
	defer ob.mu.Unlock()
	ob.buffer.Reset()
}
