package mapfile

import (
	"github.com/hetianyi/gox/file"
	"io"
	"os"
	"sync"
)

type AppendFile struct {
	bufferNum  int
	bufferSlot []byte
	bufferLock []sync.Mutex
	logSize    int
	out        *os.File
	appendFile string
	step       int
	lock       *sync.Mutex
}

func (a *AppendFile) init() error {
	a.lock.Lock()
	defer a.lock.Unlock()

	if a.bufferLock == nil {
		a.bufferLock = make([]sync.Mutex, a.bufferNum)
	}
	if a.bufferSlot == nil {
		a.bufferSlot = make([]byte, a.bufferNum)
	}
	if a.out == nil {
		if !file.Exists(a.appendFile) {
			o, err := file.CreateFile(a.appendFile)
			if err != nil {
				return err
			}
			_, err = o.WriteAt([]byte{0}, int64(a.logSize*a.bufferNum+a.bufferNum-1))
			if err != nil {
				return err
			}
			a.out = o
		} else {
			o, err := file.OpenFile(a.appendFile, os.O_RDWR, 0666)
			if err != nil {
				return err
			}

			buffers := make([]byte, a.logSize*a.bufferNum+a.bufferNum)
			if err != nil {
				return err
			}
			_, err = io.ReadAtLeast(o, buffers, len(buffers))
			if err != nil {
				return err
			}
			for i := 0; i < len(buffers); i += a.logSize + 1 {
				a.bufferSlot[i/(a.logSize+1)] = buffers[i]
			}
			a.out = o
		}
	}
	return nil
}

// TODO
func recover() {

}
