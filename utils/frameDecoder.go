package utils

import (
	"bytes"
	"sync"
)

type FrameDecoder struct {
	In      chan []byte
	Out     chan []byte
	pipline chan byte
	mu      sync.Mutex
	offset  int
	f       func([]byte) uint32
	status  bool
}

func NewFrameDecoder(LengthFieldOffset int, bytesToUint32 func([]byte) uint32) *FrameDecoder {
	return &FrameDecoder{
		In:      make(chan []byte, 256),
		Out:     make(chan []byte, 256),
		pipline: make(chan byte, 1024*64),
		offset:  LengthFieldOffset,
		f:       bytesToUint32,
		status:  false,
	}
}
func (dc *FrameDecoder) Separate() {
	dc.status = true
	go func() {
		for dc.status {
			bytes := <-dc.In
			for _, b := range bytes {
				dc.pipline <- b
			}
			bytes = nil
		}
	}()

	go func() {
		var buf bytes.Buffer
		for k := 0; dc.status; {
			if k == dc.offset {
				lb := make([]byte, 4)
				for i := 0; i < 4; i++ {
					lb[i] = <-dc.pipline
				}
				len := int(dc.f(lb))
				buf.Write(lb)
				for i := 0; i < len; i++ {
					buf.WriteByte(<-dc.pipline)
				}
				result := make([]byte, dc.offset+4+len)
				buf.Read(result)
				dc.Out <- result
				buf.Reset()
				k = 0
			} else {
				b := <-dc.pipline
				buf.WriteByte(b)
				k++
			}
		}
	}()
}
func (dc *FrameDecoder) Stop() {
	dc.status = false
}
