package nodes

import (
	"os"
	"sync"
)

var (
	Filemu sync.RWMutex
	Memmu  sync.RWMutex
)

func (m *Nodes) Save(nodesfile string) error {
	if nodesfile == "" {
		nodesfile = "./nodes"
	}
	Memmu.RLock()
	data, err := m.Marshal()
	Memmu.RUnlock()
	if err == nil {
		Filemu.Lock()
		err = os.WriteFile(nodesfile, data, 0644)
		Filemu.Unlock()
	}
	return err
}

func (m *Nodes) Load(nodesfile string) error {
	if nodesfile == "" {
		nodesfile = "./nodes"
	}
	Filemu.RLock()
	data, err := os.ReadFile(nodesfile)
	Filemu.RUnlock()
	if err == nil {
		Memmu.Lock()
		err = m.Unmarshal(data)
		Memmu.Unlock()
	} else if os.IsNotExist(err) {
		err = nil
		m.Hosts = make(map[string]uint64)
		m.Ip64S = make(map[uint64]string)
		m.Nodes = make(map[string]string)
	}
	return err
}
