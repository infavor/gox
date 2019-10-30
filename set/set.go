package set

import (
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashcode"
	"github.com/hetianyi/gox/logger"
	"sync"
)

type DataSet struct {
	m             *FixedSizeFileMap
	a             *AppendFile
	readLock      *sync.Mutex
	writeLock     *sync.Mutex
	addressBuffer []byte
}

func NewDataSet(m *FixedSizeFileMap, a *AppendFile) *DataSet {
	ret := &DataSet{
		m:             m,
		a:             a,
		readLock:      new(sync.Mutex),
		writeLock:     new(sync.Mutex),
		addressBuffer: make([]byte, 8),
	}
	return ret
}

func (d *DataSet) getIndex(data []byte) int {
	key := gox.Md5Sum(string(data))
	h := hashcode.HashCode(key)
	h = h ^ (h >> 16)
	return (d.m.slotNum - 1) & int(h)
}

func (d *DataSet) Add(data []byte) error {
	d.writeLock.Lock()
	defer d.writeLock.Unlock()

	index := d.getIndex(data)

	addr, err := d.m.Read(index)
	if err != nil {
		return err
	}
	var l int64 = 0
	if addr != nil {
		l = convert.Bytes2Length(addr)
	}
	if l == 0 {
		addr, err := d.a.ApplyAddress()
		if err != nil {
			return err
		}
		if err := d.m.Write(index, convert.Length2Bytes(addr, d.addressBuffer)); err != nil {
			return err
		}
		if err := d.a.Write(data, addr); err != nil {
			return err
		}
	} else {
		x, _, err := d.a.Contains(data, l)
		if err != nil {
			return err
		}
		if !x {
			if err := d.a.Write(data, l); err != nil {
				// logger.Info("add ", string(data))
				return err
			}
		}
	}
	return nil
}

// Remove
func (d *DataSet) Remove(data []byte) (bool, error) {
	d.writeLock.Lock()
	defer d.writeLock.Unlock()

	index := d.getIndex(data)

	addr, err := d.m.Read(index)
	if err != nil {
		return false, err
	}
	var l int64 = 0
	if addr != nil {
		l = convert.Bytes2Length(addr)
	}
	if l == 0 {
		return false, nil
	} else {
		return d.a.Delete(data, l)
	}
}

func (d *DataSet) Contains(data []byte) (bool, error) {
	addr, err := d.m.Read(d.getIndex(data))
	if err != nil {
		return false, err
	}
	var l int64 = 0
	if addr != nil {
		l = convert.Bytes2Length(addr)
	} else {
		return false, nil
	}

	x, _, err := d.a.Contains(data, l)
	if err != nil {
		return false, err
	}
	return x, nil
}
