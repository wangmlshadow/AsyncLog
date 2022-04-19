package LogModule

import (
	"runtime"
	"sync"
	"fmt"
)

const (
	logFileName = "./WebLog.txt"
)

type Impl struct {
	line int
	basename string
	stream *LogStream
}

func NewImpl(line int, basename string) *Impl {
	return &Impl{
		line: line,
		basename: basename,
		stream: NewLogStream(),
	}
}

type Logger struct {
	logFileName string
	impl *Impl
}

func NewLogger() *Logger {
	_, fileName, line, _ := runtime.Caller(2)
	fmt.Println("runtime Caller:", fileName, line)
	obj :=  &Logger{
		logFileName: logFileName,
		impl: NewImpl(line, fileName),
	}
	//runtime.SetFinalizer(obj, deleteLogger)
	return obj
}

func (logger *Logger)PrintInfo() {
	logger.impl.stream.Log("------")
	logger.impl.stream.Log(logger.impl.basename)
	logger.impl.stream.Log(" : ")
	logger.impl.stream.Log(logger.impl.line)
	logger.impl.stream.Log("\n")
	//output(logger.impl.stream.buffer.Data)
}

var once sync.Once
var asyncLogger *AsyncLogger

func once_init() {
	asyncLogger = NewAsyncLogger(logFileName, FlushInterval)
	asyncLogger.Start()
}

func output(msg []byte) {
	once.Do(once_init)
	asyncLogger.Append(msg)
}

func LOG(v interface{}) {
	//fmt.Println("LOG")
	log := NewLogger()
	log.PrintInfo()
	log.impl.stream.Log(v)
	log.impl.stream.Log("\n")
	output(log.impl.stream.buffer.Data)
}