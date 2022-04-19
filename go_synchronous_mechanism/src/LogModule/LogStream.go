package LogModule

import (
	"fmt"
)

const (
	KSmallBuffer = 4000
	KLargeBuffer = 4000 * 1000
	CurrentBuffer = 1024
)

type FixedBuffer struct {
	Data []byte // data
	//cur int
}

// func NewFixedBuffer(len int) *FixedBuffer {
// 	data := make([]byte, 0, len)
// 	return &FixedBuffer{Data: data}
// }

func NewFixedBuffer() *FixedBuffer {
	data := make([]byte, 0, CurrentBuffer)
	return &FixedBuffer{Data: data}
}

func (fixedBuffer *FixedBuffer) Append(buf []byte) {
	if (fixedBuffer.Available() > len(buf)) {
		fixedBuffer.Data = append(fixedBuffer.Data, buf...)
	} else {
		fmt.Println("fixedBuffer.Append failed")
	}
}

func (fixedBuffer *FixedBuffer) Length() int {
	return len(fixedBuffer.Data)
}

func (fixedBuffer *FixedBuffer) Reset() {
	fixedBuffer.Data = fixedBuffer.Data[:0]
}

func (fixedBuffer *FixedBuffer) Available() int {
	return cap(fixedBuffer.Data) - len(fixedBuffer.Data)
}

type LogStream struct {
	buffer *FixedBuffer
}

func NewLogStream() *LogStream {
	return &LogStream{
		buffer: NewFixedBuffer(),
	}
}

func (logStream *LogStream) Log(v interface{}) {
	s := fmt.Sprintf("%v", v)
	fmt.Println("Log:", s)
	logStream.buffer.Append([]byte(s))
}