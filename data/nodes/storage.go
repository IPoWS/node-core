package nodes

import (
	"os"
	"sync"
)

var (
	filemu sync.RWMutex
	memmu  sync.RWMutex
)

func (m *Nodes) Save(nodesfile string) error {
	if nodesfile == "" {
		nodesfile = "./nodes"
	}
	memmu.RLock()
	data, err := m.Marshal()
	memmu.RUnlock()
	if err == nil {
		filemu.Lock()
		err = os.WriteFile(nodesfile, data, 0644)
		filemu.Unlock()
	}
	return err
}

func (m *Nodes) Load(nodesfile string) error {
	if nodesfile == "" {
		nodesfile = "./nodes"
	}
	filemu.RLock()
	data, err := os.ReadFile(nodesfile)
	filemu.RUnlock()
	if err == nil {
		memmu.Lock()
		err = m.Unmarshal(data)
		memmu.Unlock()
	} else if os.IsNotExist(err) {
		err = nil
	}
	return err
}
