package router

import (
	"sync"

	"github.com/gorilla/websocket"
)

type TransItem struct {
	To         uint64
	Next       uint64
	Delay100us uint16
	Conn       *websocket.Conn
}

type transTable struct {
	// to -> ti
	table map[uint64]*TransItem
	// delay / 100us -> ti
	delays [65536]*TransItem
	mu     sync.RWMutex
}

func (t *transTable) init() {
	t.table = make(map[uint64]*TransItem)
}

func (t *transTable) add(item *TransItem) {
	t.mu.RLock()
	i, ok := t.table[item.To]
	if ok && i != nil && i.Conn != nil {
		if i.Next != item.Next {
			if item.Delay100us >= i.Delay100us {
				t.mu.RUnlock()
				return
			}
		} else {
			return
		}
	}
	t.mu.RUnlock()
	t.mu.Lock()
	t.delays[item.Delay100us] = item
	t.table[item.To] = item
	t.mu.Unlock()
}

func (t *transTable) del(to uint64) {
	t.mu.Lock()
	i, ok := t.table[to]
	if ok {
		delete(t.table, to)
		t.delays[i.Delay100us] = nil
	}
	t.mu.Unlock()
}

func (t *transTable) nextHop(to uint64) *TransItem {
	defer t.mu.RUnlock()
	t.mu.RLock()
	// 最长掩码匹配
	var i uint64 = 0xffff_ffff_ffff_ffff
	for i != 0 {
		item, ok := t.table[to&i]
		if ok {
			return item
		}
		i <<= 1
	}
	return nil
}

func (t *transTable) near() (r []*TransItem) {
	for i := 0; i < 8192; i++ {
		i := t.delays[i]
		if i != nil {
			r = append(r, i)
		}
	}
	return
}

func (t *transTable) all() (r []*TransItem) {
	t.mu.RLock()
	for _, item := range t.table {
		if item != nil {
			r = append(r, item)
		}
	}
	t.mu.RUnlock()
	return
}

func (t *transTable) isIn(ip uint64) bool {
	t.mu.RLock()
	i, ok := t.table[ip]
	t.mu.RUnlock()
	return ok && i.Conn != nil
}
