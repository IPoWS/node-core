package router

import "sync"

type transItem struct {
	to      uint64
	next    uint64
	delayms uint16
}

type TransTable struct {
	// to -> ti
	table map[uint64]*transItem
	// delay -> ti
	delays [65536]*transItem
	mu     sync.RWMutex
}

func (t *TransTable) init() {
	t.table = make(map[uint64]*transItem)
}

func (t *TransTable) add(item *transItem) {
	t.mu.Lock()
	t.delays[item.delayms] = item
	t.table[item.to] = item
	t.mu.Unlock()
}

func (t *TransTable) del(to uint64) {
	t.mu.Lock()
	i, ok := t.table[to]
	if ok {
		delete(t.table, to)
		t.delays[i.delayms] = nil
	}
	t.mu.Unlock()
}

func (t *TransTable) nextHop(to uint64) *transItem {
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
