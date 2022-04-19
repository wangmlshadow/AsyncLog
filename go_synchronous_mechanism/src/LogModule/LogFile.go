package LogModule

import (
	"sync"
)

type LogFile struct {
	basename string
	flushEveryN int
	count int
	appendFile *AppendFile
	mutex *sync.Mutex
}

func NewLogFile(basename string, flushEveryN int) *LogFile {
	appendFile := NewAppendFile(basename)
	mutex := &sync.Mutex{}
	return &LogFile{
		basename:basename,
		flushEveryN:flushEveryN,
		count:0,
		appendFile:appendFile,
		mutex:mutex,
	}
}

func (logFile *LogFile) Append(line []byte) {
	logFile.mutex.Lock()
	defer logFile.mutex.Unlock()
	logFile.append_unlock(line)
}

func (logFile *LogFile) Flush() {
	logFile.mutex.Lock()
	defer logFile.mutex.Unlock()
	logFile.appendFile.Flush()
}

func (logFile *LogFile) append_unlock(line []byte) {
	logFile.appendFile.Append(line)
	logFile.count += 1
	if logFile.count >= logFile.flushEveryN {
		logFile.appendFile.Flush()
		logFile.count = 0
	}
}