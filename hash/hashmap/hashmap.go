package hashmap

import (
	"github.com/hetianyi/gox/hash/hashcode"
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
	table []*node
	lock  *sync.Mutex
}

func NewMap() *hashMap {
	m := &hashMap{
		lock: new(sync.Mutex),
	}
	return m
}

func (m *hashMap) resize() {
	if m.table == nil {
		m.table = make([]*node, DEFAULT_INITIAL_CAPACITY)
	}
}

func (m *hashMap) Put(key interface{}, value interface{}) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.resize()
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
		for {
			if n.hashcode == h && n.value == value {
				return
			}
			n = n.next
		}
		m.table[i] = &node{
			hashcode: h,
			key:      key,
			value:    value,
			next:     m.table[0],
		}
	}
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
	for {
		if n.hashcode == h {
			return n.value
		}
		n = n.next
	}
}

func (m *hashMap) getIndex(key interface{}) (int, int32) {
	h := hashcode.HashCode(key)
	return int((len(m.table) - 1) & int(h)), h
}
