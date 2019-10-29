package set

import (
	"bytes"
	"errors"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/file"
	"os"
	"sync"
)

var (
	zeroByte = []byte{0}
)

type AppendFile struct {
	tailAddr     [9]byte
	bufferLock   *sync.Mutex
	logSize      int
	out          *os.File
	in           *os.File
	appendFile   string
	step         int   // continuous space for every slot
	curOffset    int64 // current write offset
	writeLock    *sync.Mutex
	applyLock    *sync.Mutex
	oneByteArray []byte
	buffer       *bytes.Buffer
	stepBuff     []byte
}

func NewAppendFile(logSize, step int, appendFile string) (*AppendFile, error) {
	r := &AppendFile{
		tailAddr:     [9]byte{},
		bufferLock:   new(sync.Mutex),
		writeLock:    new(sync.Mutex),
		applyLock:    new(sync.Mutex),
		oneByteArray: make([]byte, 1),
		logSize:      logSize,
		step:         step,
		appendFile:   appendFile,
		buffer:       new(bytes.Buffer),
		stepBuff:     make([]byte, (logSize+1)*step+9),
	}
	return r, r.init()
}

func (a *AppendFile) init() (err error) {
	a.writeLock.Lock()
	defer func() {
		if err != nil && a.out != nil {
			a.out.Close()
		}
		a.writeLock.Unlock()
	}()

	if a.out == nil {
		if !file.Exists(a.appendFile) {
			o, err := file.CreateFile(a.appendFile)
			if err != nil {
				return err
			}
			a.out = o
		} else {
			o, err := file.OpenFile(a.appendFile, os.O_RDWR, 0666)
			if err != nil {
				return err
			}
			a.out = o
		}
	}
	if a.in == nil {
		i, err := file.OpenFile(a.appendFile, os.O_RDONLY, 0666)
		if err != nil {
			return err
		}
		a.in = i
	}
	fInfo, err := a.in.Stat()
	if err != nil {
		return err
	}
	if fInfo.Size() == 0 {
		if err := a.extend(nil); err != nil {
			return err
		}
	} else {
		a.curOffset = fInfo.Size()
	}
	return nil
}

func (a *AppendFile) ApplyAddress() (int64, error) {
	a.applyLock.Lock()
	defer a.applyLock.Unlock()

	if err := a.extend(nil); err != nil {
		return -1, err
	}
	return a.curOffset - int64((a.logSize+1)*a.step+9), nil
}

func (a *AppendFile) Write(data []byte, offset int64) error {
	a.writeLock.Lock()
	defer func() {
		a.buffer.Reset()
		a.writeLock.Unlock()
	}()

	if len(data) != a.logSize {
		return errors.New("data does not match log size")
	}

	a.buffer.Write(data)
	a.buffer.WriteByte(1)

	return a.append(offset)
}

func (a *AppendFile) Delete(data []byte, offset int64) (bool, error) {
	a.writeLock.Lock()
	defer func() {
		a.buffer.Reset()
		a.writeLock.Unlock()
	}()

	if len(data) != a.logSize {
		return false, errors.New("data does not match log size")
	}
	a.buffer.Write(data)
	a.buffer.WriteByte(1)
	return a.delete(offset)
}

func (a *AppendFile) Contains(data []byte, offset int64) (bool, int, error) {
	return a.read(data, offset, 0)
}

func (a *AppendFile) read(data []byte, blockHeadOffset int64, depth int) (bool, int, error) {
	stepBuff := make([]byte, (a.logSize+1)*a.step+9)
	if _, err := a.out.ReadAt(stepBuff, blockHeadOffset); err != nil {
		return false, 0, err
	}
	for i := 0; i < a.step; i++ {
		depth++
		if stepBuff[a.logSize*(i+1)+i] == 1 &&
			bytes.Equal(data, stepBuff[(a.logSize+1)*i:(a.logSize+1)*i+a.logSize]) {
			return true, depth, nil
		}
	}
	if stepBuff[len(stepBuff)-1] == 0 {
		return false, 0, nil
	}
	return a.read(data, convert.Bytes2Length(stepBuff[len(stepBuff)-9:len(stepBuff)-1]), depth)
}

func copy(target []byte, src []byte, offset int) {
	for i := 0; i < len(src); i++ {
		target[offset+i] = src[i]
	}
}

func (a *AppendFile) extend(baseData []byte) error {
	t := make([]byte, (a.logSize+1)*a.step+9)
	copy(t, baseData, len(baseData))
	if _, err := a.out.Seek(0, 2); err != nil {
		return err
	}
	if _, err := a.out.Write(t); err != nil {
		return err
	}
	a.curOffset += int64((a.logSize+1)*a.step) + 9
	return nil
}

func (a *AppendFile) append(blockHeadOffset int64) error {
	// read placeholder
	for i := 0; i < a.step; i++ {
		t, err := a.readOneByte(blockHeadOffset + int64(a.logSize*(i+1)+i))
		if err != nil {
			return err
		}
		// already has data
		if t[0] == 1 {
			continue
		}
		if _, err := a.out.WriteAt(a.buffer.Bytes(),
			blockHeadOffset+int64(a.logSize*i+i)); err != nil {
			return err
		}
		return nil
	}
	// get next block address
	if _, err := a.out.ReadAt(a.tailAddr[:],
		blockHeadOffset+int64((a.logSize+1)*a.step)); err != nil {
		return err
	}
	// valid address data, continue next block.
	if a.tailAddr[8] == 1 {
		nextBlockHeadOffset := convert.Bytes2Length(a.tailAddr[0:8])
		return a.append(nextBlockHeadOffset)
	}
	addr, err := a.ApplyAddress()
	if err != nil {
		return err
	}
	// write address
	convert.Length2Bytes(addr, a.tailAddr[0:8])
	a.tailAddr[8] = 1
	if _, err = a.out.WriteAt(a.tailAddr[:], blockHeadOffset+int64((a.logSize+1)*a.step)); err != nil {
		return err
	}
	return a.append(addr)
}

func (a *AppendFile) delete(blockHeadOffset int64) (bool, error) {
	if _, err := a.out.ReadAt(a.stepBuff, blockHeadOffset); err != nil {
		return false, err
	}
	for i := 0; i < a.step; i++ {
		if a.stepBuff[a.logSize*(i+1)+i] == 1 &&
			bytes.Equal(a.buffer.Bytes(), a.stepBuff[(a.logSize+1)*i:(a.logSize+1)*(i+1)]) {
			if _, err := a.out.WriteAt(empty, blockHeadOffset+int64(a.logSize*(i+1)+i)); err != nil {
				return false, err
			}
			return true, nil
		}
	}
	// valid address data, continue next block.
	if a.stepBuff[len(a.stepBuff)-1] == 1 {
		nextBlockHeadOffset := convert.Bytes2Length(a.stepBuff[len(a.stepBuff)-9 : len(a.stepBuff)-1])
		return a.delete(nextBlockHeadOffset)
	}
	return false, nil
}

func (a *AppendFile) readOneByte(offset int64) ([]byte, error) {
	t := make([]byte, 1)
	_, err := a.out.ReadAt(t, offset)
	if err != nil {
		return nil, err
	}
	return t, nil
}
