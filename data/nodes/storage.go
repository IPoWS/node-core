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
	if err == nil && data != nil && len(data) > 0 {
		Memmu.Lock()
		err = m.Unmarshal(data)
		Memmu.Unlock()
	} else {
		m.Hosts = make(map[string]uint64)
		m.Ip64S = make(map[uint64]string)
		m.Nodes = make(map[string]string)
		m.Times = make(map[uint64]uint64)
	}
	return err
}

func (m *Nodes) CopyNodes() map[string]string {
	Memmu.RLock()
	ret := make(map[string]string, len(m.Nodes))
	for k, v := range m.Nodes {
		ret[k] = v
	}
	Memmu.RUnlock()
	return ret
}

func (m *Nodes) CopyIp64S() map[uint64]string {
	Memmu.RLock()
	ret := make(map[uint64]string, len(m.Ip64S))
	for k, v := range m.Ip64S {
		ret[k] = v
	}
	Memmu.RUnlock()
	return ret
}

func (m *Nodes) CopyHosts() map[string]uint64 {
	Memmu.RLock()
	ret := make(map[string]uint64, len(m.Hosts))
	for k, v := range m.Hosts {
		ret[k] = v
	}
	Memmu.RUnlock()
	return ret
}

func (m *Nodes) CopyTimes() map[uint64]uint64 {
	Memmu.RLock()
	ret := make(map[uint64]uint64, len(m.Times))
	for k, v := range m.Times {
		ret[k] = v
	}
	Memmu.RUnlock()
	return ret
}
