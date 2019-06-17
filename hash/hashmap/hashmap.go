// hashmap比golang内置的map添加速度要快，1000w数据hashmap比内置map快几百毫秒-1秒
// 删除速度也较快
// 内存占用由于hashmap实现的原因，在接近2^n但不到2^n时，hashmap占用空间较小，当hashmap扩容时（此时元素刚刚超过2^n，但是容量翻翻），hashmap会比内置map占用内存大。
package hashmap

import (
	"errors"
	"github.com/hetianyi/gox"
	"github.com/hetianyi/gox/convert"
	"github.com/hetianyi/gox/hash/hashcode"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
)

const (
	DEFAULT_INITIAL_CAPACITY = 1 << 4
	MAXIMUM_CAPACITY         = 1 << 30 // 最大容量
	DEFAULT_LOAD_FACTOR      = 0.75    // 阈值比率
)

type node struct {
	hashcode int32
	key      interface{}
	value    interface{}
	next     *node
}

type hashMap struct {
	table      []*node
	size       int
	threshold  int // 阈值
	lock       *sync.Mutex
	loadFactor float64
}

func NewMap() *hashMap {
	m := &hashMap{
		lock:       new(sync.Mutex),
		loadFactor: DEFAULT_LOAD_FACTOR,
	}
	return m
}

func New(initialCapacity int, loadFactor float64) *hashMap {
	if loadFactor <= 0 || math.IsNaN(loadFactor) {
		panic(errors.New("Illegal load factor: " + convert.Float64ToStr(loadFactor)))
	}
	if initialCapacity > MAXIMUM_CAPACITY {
		initialCapacity = MAXIMUM_CAPACITY
	}

	m := &hashMap{
		lock:       new(sync.Mutex),
		loadFactor: loadFactor,
		threshold:  tableSizeFor(initialCapacity),
	}
	return m
}

func (m *hashMap) resize() {
	oldCap := gox.TValue(m.table == nil || len(m.table) == 0, 0, len(m.table)).(int)
	oldThr := m.threshold
	var newCap, newThr int

	if oldCap > 0 {
		if oldCap >= MAXIMUM_CAPACITY {
			m.threshold = math.MaxInt32
			log.Trace("resize map finish")
			return
		}
		if newCap = oldCap << 1; newCap < MAXIMUM_CAPACITY && oldCap >= DEFAULT_INITIAL_CAPACITY {
			newThr = oldThr << 1 // double threshold
		}
	} else if oldThr > 0 {
		newCap = oldThr
	} else {
		newCap = DEFAULT_INITIAL_CAPACITY
		newThr = (int)(DEFAULT_LOAD_FACTOR * DEFAULT_INITIAL_CAPACITY)
	}
	if newThr == 0 {
		ft := float64(newCap) * m.loadFactor
		newThr = gox.TValue(newCap < MAXIMUM_CAPACITY && ft < float64(MAXIMUM_CAPACITY), int(ft), math.MaxInt32).(int)
	}
	m.threshold = newThr

	log.Trace("resize map, current capacity is ", len(m.table), ", new capacity is ", newCap, ", old size is ", m.size)
	newTab := make([]*node, newCap)
	if m.table != nil {
		for i := 0; i < len(m.table); i++ {
			e := m.table[i]
			if e != nil {
				m.table[i] = nil
				if e.next == nil {
					newTab[int(int(e.hashcode)&(newCap-1))] = e
				} else {
					for e != nil {
						cp := e
						e = e.next
						cp.next = nil
						idx := int(int(cp.hashcode) & (newCap - 1))
						newE := newTab[idx]
						if newE == nil {
							newTab[idx] = cp
						} else {
							for newE != nil {
								if newE.next == nil {
									newE.next = cp
									break
								}
								newE = newE.next
							}
						}
					}
				}
			}
		}
	}
	m.table = newTab
	log.Trace("resize map finish")
}

func (m *hashMap) Put(key interface{}, value interface{}) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	if m.table == nil || len(m.table) == 0 {
		m.resize()
	}
	i, h := m.getIndex(key)
	if m.table[i] == nil {
		m.table[i] = &node{
			hashcode: h,
			key:      key,
			value:    value,
			next:     nil,
		}
	} else {
		n := m.table[i]
		for n != nil {
			if n.hashcode == h && n.value == value {
				log.Trace("replace old value ", n.value, " to ", value)
				oldVal := n.value
				n.value = value
				return oldVal
			}
			n = n.next
		}
		m.table[i] = &node{
			hashcode: h,
			key:      key,
			value:    value,
			next:     m.table[i],
		}
	}
	m.size++
	if m.size > m.threshold {
		m.resize()
	}
	return nil
}

func (m *hashMap) Get(key interface{}) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	i, h := m.getIndex(key)
	n := m.table[i]
	if n == nil {
		return nil
	}
	for n != nil {
		if n.hashcode == h && n.key == key {
			return n.value
		}
		n = n.next
	}
	return nil
}

func (m *hashMap) getIndex(key interface{}) (int, int32) {
	h := hashcode.HashCode(key)
	return int((len(m.table) - 1) & int(h)), h
}

func (m *hashMap) Remove(key interface{}) interface{} {
	m.lock.Lock()
	defer m.lock.Unlock()
	i, h := m.getIndex(key)
	n := m.table[i]
	if n == nil {
		return nil
	}
	var prev *node = nil
	for n != nil {
		if n.hashcode == h && n.key == key {
			ret := n.value
			if prev != nil {
				prev.next = n.next
			}
			return ret
		}
		prev = n
		n = n.next
	}
	return nil
}

func tableSizeFor(cap int) int {
	n := cap - 1
	n |= n >> 1
	n |= n >> 2
	n |= n >> 4
	n |= n >> 8
	n |= n >> 16
	return gox.TValue(n < 0, 1, gox.TValue(n >= 1073741824, 1073741824, n+1).(int)).(int)
}
