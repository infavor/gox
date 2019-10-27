package set

import (
	"errors"
	"github.com/hetianyi/gox/file"
	"os"
	"sync"
)

var valued = []byte{1}
var empty = []byte{0}

// FixedSizeFileMap is a fixed size file map.
type FixedSizeFileMap struct {
	slotNum       int      // number of slots
	slotSize      int      // byte size of each slot
	out           *os.File // target binlog file
	binlogFile    string   // target binlog file
	slotMap       []byte   // binlog slot map, this is stored in memory
	slotLockMap   map[int]*sync.Mutex
	lock          *sync.Mutex
	slotWriteLock *sync.Mutex
}

// NewFileMap creates a new FixedSizeFileMap.
func NewFileMap(slotNum, slotSize int, binlogFile string) (*FixedSizeFileMap, error) {
	m := &FixedSizeFileMap{
		slotNum:       slotNum,
		slotSize:      slotSize,
		binlogFile:    binlogFile,
		slotLockMap:   make(map[int]*sync.Mutex),
		lock:          new(sync.Mutex),
		slotWriteLock: new(sync.Mutex),
		slotMap:       make([]byte, slotNum),
	}
	return m, m.init()
}

func (m *FixedSizeFileMap) init() error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if m.slotMap == nil {
		m.slotMap = make([]byte, m.slotNum)
	}
	if m.out == nil {
		if !file.Exists(m.binlogFile) {
			o, err := file.CreateFile(m.binlogFile)
			if err != nil {
				return err
			}
			_, err = o.WriteAt([]byte{0}, int64(m.slotNum*m.slotSize+m.slotNum-1))
			if err != nil {
				return err
			}
			m.out = o
		} else {
			o, err := file.OpenFile(m.binlogFile, os.O_RDWR, 0666)
			if err != nil {
				return err
			}
			_, err = o.ReadAt(m.slotMap, 0)
			if err != nil {
				return err
			}
			m.out = o
		}
	}
	return nil
}

func (m *FixedSizeFileMap) SlotSnapshot() []byte {
	m.lock.Lock()
	defer m.lock.Unlock()

	s := make([]byte, len(m.slotMap))
	copy(s, m.slotMap, 0)
	return s
}

func (m *FixedSizeFileMap) lockSlot(slotIndex int) *sync.Mutex {
	m.slotWriteLock.Lock()
	defer m.slotWriteLock.Unlock()
	lo := m.slotLockMap[slotIndex]
	if lo == nil {
		lo = new(sync.Mutex)
		m.slotLockMap[slotIndex] = lo
	}
	return lo
}

// Write writes data in a slot.
//
//  slotIndex begin from 0,
//  data is slot data.
func (m *FixedSizeFileMap) Write(slotIndex int, data []byte) error {
	m.lock.Lock()
	lo := m.lockSlot(slotIndex)
	if lo != nil {
		lo.Lock()
	}
	defer func() {
		if lo != nil {
			delete(m.slotLockMap, slotIndex)
			lo.Unlock()
		}
		m.lock.Unlock()
	}()

	if slotIndex < 0 || slotIndex >= m.slotNum {
		return errors.New("index of out range")
	}
	if m.slotMap[slotIndex] == 1 {
		return errors.New("write failed: slot already has data")
	}
	if len(data) != m.slotSize {
		return errors.New("data size mismatch the slot size")
	}

	_, err := m.out.WriteAt(data, int64(m.slotNum)+int64((slotIndex)*m.slotSize))
	if err != nil {
		return err
	}
	_, err = m.out.WriteAt(valued, int64(slotIndex))
	if err != nil {
		return err
	}
	m.slotMap[slotIndex] = 1
	return nil
}

// Delete deletes data of a slot.
//
//  slotIndex begin from 0,
func (m *FixedSizeFileMap) Delete(slotIndex int) error {
	m.lock.Lock()
	lo := m.lockSlot(slotIndex)
	if lo != nil {
		lo.Lock()
	}
	defer func() {
		if lo != nil {
			delete(m.slotLockMap, slotIndex)
			lo.Unlock()
		}
		m.lock.Unlock()
	}()

	if slotIndex < 0 || slotIndex >= m.slotNum {
		return errors.New("index of out range")
	}
	if m.slotMap[slotIndex] == 0 {
		return nil
	}

	_, err := m.out.WriteAt(empty, int64(slotIndex))
	if err != nil {
		return err
	}
	m.slotMap[slotIndex] = 0
	return nil
}

// Read reads slot data from binlog file.
func (m *FixedSizeFileMap) Read(slotIndex int) ([]byte, error) {
	lo := m.lockSlot(slotIndex)
	if lo != nil {
		lo.Lock()
	}
	defer func() {
		if lo != nil {
			lo.Unlock()
		}
	}()

	if slotIndex < 0 || slotIndex >= m.slotNum {
		return nil, errors.New("index of out range")
	}
	if m.slotMap[slotIndex] == 0 {
		return nil, nil
	}
	ret := make([]byte, m.slotSize)
	_, err := m.out.ReadAt(ret, int64(m.slotNum)+int64((slotIndex)*m.slotSize))
	return ret, err
}
