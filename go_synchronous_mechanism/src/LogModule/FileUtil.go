package LogModule

import (
	"os"
	"fmt"
	//"bufio"
	"runtime"
)

const (
	BUFSIZE = 1024
)

type AppendFile struct {
	fp *os.File
	//buf *bufio.Writer // buffer
	buf []byte //
}

func deleteAppendFile(appendFile *AppendFile) {
	fmt.Println("deleteAppendFile")
	appendFile.fp.Close()
}

func NewAppendFile(filename string) (*AppendFile) {
	fp, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		//fmt.Println("OpenFile error: ", err)
		//return nil, err
		panic(err)
	}
	//buf := bufio.NewWriterSize(fp, 64*1024)
	buf := make([]byte, 0, BUFSIZE)
	obj := &AppendFile{fp:fp, buf:buf}
	runtime.SetFinalizer(obj, deleteAppendFile)
	return obj
}

func (appendFile *AppendFile) Flush() error {
	fmt.Println("appendFile flush")
	//return appendFile.buf.Flush()
	_, err := appendFile.fp.Write(appendFile.buf)
	appendFile.buf = appendFile.buf[:0]
	if err != nil {
		//fmt.Println("Write error: ", err)
		panic(err)
	}
	return appendFile.fp.Sync()
}

func (appendFile *AppendFile) Append(b []byte) {
	//_, err := appendFile.buf.Write(b)
	if len(b) > cap(appendFile.buf) - len(appendFile.buf) {
		appendFile.Flush()
	}
	if len(b) > cap(appendFile.buf) - len(appendFile.buf) {
		appendFile.fp.Write(b)
	} else {
		appendFile.buf = append(appendFile.buf, b...)
	}
}