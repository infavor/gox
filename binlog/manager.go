package binlog

import (
	"errors"
	"github.com/hetianyi/gox/file"
	"os"
	"sync"
)

// FixedSizeBinlogManager is a fixed size binlog file.
type FixedSizeBinlogManager struct {
	slotNum    int      // number of slots
	slotSize   int      // byte size of each slot
	out        *os.File // target binlog file
	binlogFile string   // target binlog file
	slotMap    []byte   // binlog slot map, this is stored in memory

	lock *sync.Mutex
}

func (m *FixedSizeBinlogManager) init() error {
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
			m.out = o
		}
	}
	return nil
}

// Write writes data in a slot.
//
//  slotIndex begin from 0,
//  data is slot data.
func (m *FixedSizeBinlogManager) Write(slotIndex int, data []byte) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	if slotIndex < 0 || slotIndex >= m.slotNum {
		return errors.New("index of out range")
	}
	if m.slotMap[slotIndex] == 1 {
		return errors.New("write failed: slot already has data")
	}
	if len(data) != m.slotSize {
		return errors.New("data size mismatch the slot size")
	}

	if data == nil {
		_, err := m.out.WriteAt([]byte{0}, int64(slotIndex))
		if err != nil {
			return err
		}
	} else {
		_, err := m.out.WriteAt(data, int64(m.slotNum)+int64((slotIndex)*m.slotSize))
		if err != nil {
			return err
		}
		_, err = m.out.WriteAt([]byte{1}, int64(slotIndex))
		if err != nil {
			return err
		}
	}
	return nil
}

// Read reads slot data from binlog file.
func (m *FixedSizeBinlogManager) Read(slotIndex int) ([]byte, error) {
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
