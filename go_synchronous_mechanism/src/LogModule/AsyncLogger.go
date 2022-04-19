package LogModule

import (
	"sync"
	"runtime"
	//"time"
	"fmt"
)

const (
	FlushInterval = 1024
)

type AsyncLogger struct {
	running bool
	basename string
	mutex *sync.Mutex
	cond *sync.Cond
	wg *sync.WaitGroup
	latch *CountDownLatch
	currentBuffer *FixedBuffer
	nextBuffer *FixedBuffer
	buffers []*FixedBuffer
}

func NewAsyncLogger(basename string, flushInterval int) *AsyncLogger {
	mutex := &sync.Mutex{}
	cond := sync.NewCond(mutex)
	obj := &AsyncLogger{
		running: false,
		basename: basename,
		mutex: mutex,
		cond: cond,
		wg: &sync.WaitGroup{},
		latch: NewCountDownLatch(1),
		currentBuffer: NewFixedBuffer(),
		nextBuffer: NewFixedBuffer(),
		buffers: make([]*FixedBuffer, 0, 16),
	}
	runtime.SetFinalizer(obj, deleteAsycLogger)
	return obj
}

func deleteAsycLogger(asyncLogger *AsyncLogger) {
	fmt.Println("delete asynclogger")
	if asyncLogger.running {
		fmt.Println("stop")
		asyncLogger.Stop()
	}
}

func (asyncLogger *AsyncLogger) Start() {
	asyncLogger.running = true
	asyncLogger.wg.Add(1)
	go asyncLogger.run()
	fmt.Println("latch wait")
	asyncLogger.latch.Wait()
	fmt.Println("latch break wait and start end")
}

func (asyncLogger *AsyncLogger) Stop() {
	fmt.Println("asynclogger stop")
	asyncLogger.running = false
	asyncLogger.cond.Signal()
	asyncLogger.wg.Wait()
}

func (asyncLogger *AsyncLogger) Append(logline []byte) {
	fmt.Println("append...")
	asyncLogger.mutex.Lock()
	defer asyncLogger.mutex.Unlock()
	if asyncLogger.currentBuffer.Available() > len(logline) {
		fmt.Println("append.....avaliabel")
		asyncLogger.currentBuffer.Append(logline)
		//asyncLogger.cond.Signal()
	} else {
		fmt.Println("append.....unavaliabel")
		asyncLogger.buffers = append(asyncLogger.buffers, asyncLogger.currentBuffer)
		asyncLogger.currentBuffer = nil
		if asyncLogger.nextBuffer != nil {
			asyncLogger.currentBuffer = asyncLogger.nextBuffer
			asyncLogger.nextBuffer = nil
		} else {
			asyncLogger.currentBuffer = NewFixedBuffer()
		}
		asyncLogger.currentBuffer.Append(logline)
		asyncLogger.cond.Signal()
		fmt.Println("buffers size:", len(asyncLogger.buffers))
	}
}

func (asyncLogger *AsyncLogger) run() {
	fmt.Println("run...")
	if !asyncLogger.running {
		panic("asyncLogger.running is not true when exec run")
	}
	asyncLogger.latch.CountDown()
	logFile := NewLogFile(asyncLogger.basename, FlushInterval)
	newbuffer1 := NewFixedBuffer()
	newbuffer2 := NewFixedBuffer()
	//buffers2write := make([]*FixedBuffer, 0, 16)
	var buffers2write []*FixedBuffer
	for asyncLogger.running {

		{
			asyncLogger.mutex.Lock()
			defer asyncLogger.mutex.Unlock()
			if len(asyncLogger.buffers) == 0 {
				// 不咋用
				fmt.Println("cond wait...")
				asyncLogger.cond.Wait()
			}
			fmt.Println("cond break wait...")
			asyncLogger.buffers = append(asyncLogger.buffers, asyncLogger.currentBuffer)
			asyncLogger.currentBuffer = newbuffer1
			buffers2write = asyncLogger.buffers
			asyncLogger.buffers = make([]*FixedBuffer, 0, 16)
			if asyncLogger.nextBuffer == nil {
				asyncLogger.nextBuffer = newbuffer2
			}
		}

		if len(buffers2write) == 0 {
			panic("buffer2write is empty")
		}

		if len(buffers2write) > 25 {
			buffers2write = buffers2write[:2]
		}

		for _, buffer := range buffers2write {
			logFile.Append(buffer.Data)
		}

		if (len(buffers2write) > 2) {
			buffers2write = buffers2write[:2]
		}

		if newbuffer1 == nil {
			if len(buffers2write) > 0 {
				newbuffer1 = buffers2write[0]
				buffers2write = buffers2write[1:]
			}
		}

		if newbuffer2 == nil {
			if len(buffers2write) > 0 {
				newbuffer2 = buffers2write[0]
				buffers2write = buffers2write[1:]
			}
		}

		buffers2write = buffers2write[:0]
		logFile.Flush()
		fmt.Println("run flush 1")
	}
	logFile.Flush()
	fmt.Println("run flush 2")
	asyncLogger.wg.Done()
}