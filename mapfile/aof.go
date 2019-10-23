package mapfile

import (
	"bytes"
	"errors"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"io"
	"os"
	"sync"
)

var (
	zeroByte = []byte{0}
)

type AppendFile struct {
	tailAddr     [9]byte
	bufferSlot   []byte
	bufferLock   *sync.Mutex
	logSize      int
	out          *os.File
	appendFile   string
	step         int   // continuous space for every slot
	curOffset    int64 // current write offset
	lock         *sync.Mutex
	oneByteArray []byte
	buffer       *bytes.Buffer
}

func NewAppendFile(logSize, step int, appendFile string) (*AppendFile, error) {
	r := &AppendFile{
		tailAddr:     [9]byte{},
		bufferSlot:   make([]byte, logSize+9),
		bufferLock:   new(sync.Mutex),
		lock:         new(sync.Mutex),
		oneByteArray: make([]byte, 1),
		logSize:      logSize,
		step:         step,
		appendFile:   appendFile,
		buffer:       new(bytes.Buffer),
	}
	return r, r.init()
}

func (a *AppendFile) init() (err error) {
	a.lock.Lock()
	defer func() {
		if a.out != nil {
			a.out.Close()
		}
		a.lock.Unlock()
	}()

	bufSize := a.logSize + 9

	if a.bufferLock == nil {
		a.bufferLock = new(sync.Mutex)
	}
	if a.bufferSlot == nil {
		a.bufferSlot = make([]byte, bufSize) // logSize + len(tailAddr) + 1
	}
	if a.out == nil {
		if !file.Exists(a.appendFile) {
			o, err := file.CreateFile(a.appendFile)
			if err != nil {
				return err
			}
			_, err = o.WriteAt([]byte{0}, int64(bufSize))
			if err != nil {
				return err
			}
			if err != nil {
				return err
			}
			a.out = o
		} else {
			o, err := file.OpenFile(a.appendFile, os.O_RDWR, 0666)
			if err != nil {
				return err
			}

			_, err = io.ReadAtLeast(o, a.bufferSlot, bufSize)
			if err != nil {
				return err
			}
			a.recover()
			a.out = o
		}
	}
	fInfo, err := a.out.Stat()
	if err != nil {
		return err
	}
	a.curOffset = fInfo.Size()
	return nil
}

func (a *AppendFile) Write(data []byte, offset int64) (int64, error) {
	a.lock.Lock()
	defer func() {
		a.bufferSlot[len(a.bufferSlot)-1] = 0
		a.buffer.Reset()
		a.lock.Unlock()
	}()

	if len(data) != a.logSize {
		return -1, errors.New("data does not match log size")
	}

	a.buffer.Write(data)
	a.buffer.Write(convert.Length2Bytes(offset, a.tailAddr[0:8]))
	a.buffer.WriteByte(1)
	copy(a.bufferSlot, a.buffer.Bytes(), 0)

	if _, err := a.out.WriteAt(a.bufferSlot, 0); err != nil {
		return -1, err
	}
	return a.append(offset)
}

func copy(target []byte, src []byte, offset int) {
	for i := 0; i < len(src); i++ {
		target[offset+i] = src[i]
	}
}

func (a *AppendFile) extend(basedata []byte) error {
	t := make([]byte, (a.logSize+1)*a.step+9)
	copy(t, basedata, len(basedata))
	if _, err := a.out.Write(t); err != nil {
		return err
	}
	a.curOffset += int64((a.logSize+1)*a.step) + 9
	return nil
}

func (a *AppendFile) append(blockHeadOffset int64) (int64, error) {
	// read placeholder
	for i := 0; i < a.step; i++ {
		t, err := a.readOneByte(blockHeadOffset + int64(a.logSize*(i+1)+i))
		if err != nil {
			return -1, err
		}
		// already has data
		if t[0] == 1 {
			continue
		}
		a.bufferSlot[(len(a.bufferSlot) - 9)] = 1
		if _, err := a.out.WriteAt(a.bufferSlot[0:(len(a.bufferSlot)-8)],
			blockHeadOffset+int64(a.logSize*i+i)); err != nil {
			return -1, err
		}
		// reset cache
		if _, err := a.out.WriteAt(zeroByte, int64(len(a.bufferSlot)-1)); err != nil {
			return -1, err
		}
		a.bufferSlot[len(a.bufferSlot)-1] = 0
		return blockHeadOffset, nil
	}
	// get next block address
	if _, err := a.out.ReadAt(a.tailAddr[:], blockHeadOffset+int64((a.logSize+1)*a.step)); err != nil {
		return -1, err
	}
	// valid address data, continue next block.
	if a.tailAddr[8] == 1 {
		nextBlockHeadOffset := convert.Bytes2Length(a.tailAddr[0:8])
		return a.append(nextBlockHeadOffset)
	}

	// block has no space , extends file.
	a.bufferSlot[(len(a.bufferSlot) - 9)] = 1
	tail := make([]byte, a.logSize*a.step+a.step+9-a.logSize-1)
	data := append(a.bufferSlot[0:(len(a.bufferSlot)-8)], tail...)
	if _, err := a.out.Write(data); err != nil {
		return -1, err
	}
	// write next address pointer.
	if _, err := a.out.WriteAt(
		append(convert.Length2Bytes(a.curOffset+int64(len(data)), a.tailAddr[:]), 1),
		blockHeadOffset+int64((a.logSize+1)*a.step)); err != nil {
		return -1, err
	}
	a.curOffset += int64((a.logSize+1)*a.step) + 9

	// reset cache
	if _, err := a.out.WriteAt(zeroByte, int64(len(a.bufferSlot)-1)); err != nil {
		return -1, err
	}
	a.bufferSlot[len(a.bufferSlot)-1] = 0
	return a.curOffset, nil
}

func (a *AppendFile) recover() (int64, error) {
	if a.bufferSlot[(len(a.bufferSlot)-1)] == 1 {
		offset := convert.Bytes2Length(a.bufferSlot[len(a.bufferSlot)-9 : len(a.bufferSlot)-1])
		if offset < int64(len(a.bufferSlot)) || offset > a.curOffset {
			return -1, errors.New("invalid offset value: " + convert.Int64ToStr(offset))
		}
		return a.append(offset)
	}
	return -1, nil
}

func (a *AppendFile) readOneByte(offset int64) ([]byte, error) {
	t := make([]byte, 1)
	_, err := a.out.ReadAt(t, offset)
	if err != nil {
		return nil, err
	}
	return t, nil
}
