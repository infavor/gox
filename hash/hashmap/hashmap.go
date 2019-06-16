package hashmap

import (
	"github.com/hetianyi/gox"
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
	table     []*node
	size      int
	threshold int // 阈值
	lock      *sync.Mutex
}

func NewMap() *hashMap {
	m := &hashMap{
		lock: new(sync.Mutex),
	}
	return m
}

func (m *hashMap) resize() {
	log.Info("resize map, current size is ", len(m.table))
	oldCap := gox.TValue(m.table == nil || len(m.table) == 0, 0, len(m.table)).(int)
	oldThr := m.threshold
	var newCap, newThr int

	if oldCap > 0 {
		if oldCap >= MAXIMUM_CAPACITY {
			m.threshold = math.MaxInt32
			log.Info("resize map finish")
			return
		} else if oldCap<<1 < MAXIMUM_CAPACITY && oldCap >= DEFAULT_INITIAL_CAPACITY {
			newCap = oldCap << 1
			newThr = oldThr << 1 // double threshold
		}
	} else if oldThr > 0 {
		newCap = oldThr
	} else {
		newCap = DEFAULT_INITIAL_CAPACITY
		newThr = (int)(DEFAULT_LOAD_FACTOR * DEFAULT_INITIAL_CAPACITY)
	}
	if newThr == 0 {
		ft := (float64)(newCap * DEFAULT_INITIAL_CAPACITY)
		newThr = gox.TValue(newCap < MAXIMUM_CAPACITY && ft < float64(MAXIMUM_CAPACITY), int(ft), math.MaxInt32).(int)
	}
	m.threshold = newThr

	newTab := make([]*node, newCap)
	if m.table != nil {
		for i := 0; i < m.size; i++ {
			e := m.table[i]
			if e != nil {
				if e.next == nil {
					newTab[int(int(e.hashcode)&(newCap-1))] = e
					m.table[i] = nil
				} else {
					for e != nil {
						//fmt.Println("2")
						cp := &node{
							hashcode: e.hashcode,
							key:      e.key,
							value:    e.value,
							next:     nil,
						}
						idx := int(int(cp.hashcode) & (newCap - 1))
						newE := newTab[idx]
						if newE == nil {
							newTab[idx] = cp
						} else {
							next := newE
							for next != nil {
								//fmt.Println("1", next.key)
								if next.next == nil {
									next.next = cp
									break
								}
								next = next.next
							}
						}
						e = e.next
					}
				}
			}
		}
	}
	m.table = newTab
	log.Info("resize map finish")
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
	m.resize()
	i, h := m.getIndex(key)
	n := m.table[i]
	if n == nil {
		return nil
	}
	for n != nil {
		if n.hashcode == h {
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
