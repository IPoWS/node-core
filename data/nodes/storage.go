package nodes

import (
	"os"
	"sync"
)

var (
	Filemu    sync.RWMutex
	Memmu     sync.RWMutex
	Nodesfile string
)

func (m *Nodes) Save() error {
	if Nodesfile == "" {
		Nodesfile = "./nodes"
	}
	Memmu.RLock()
	data, err := m.Marshal()
	Memmu.RUnlock()
	if err == nil {
		Filemu.Lock()
		err = os.WriteFile(Nodesfile, data, 0644)
		Filemu.Unlock()
	}
	return err
}

func (m *Nodes) Load() error {
	if Nodesfile == "" {
		Nodesfile = "./nodes"
	}
	Filemu.RLock()
	data, err := os.ReadFile(Nodesfile)
	Filemu.RUnlock()
	if err == nil {
		Memmu.Lock()
		err = m.Unmarshal(data)
		Memmu.Unlock()
	} else if os.IsNotExist(err) {
		Memmu.Lock()
		m.Nodes = make(map[string]*Node)
		Memmu.Unlock()
		err = nil
	}
	return err
}
