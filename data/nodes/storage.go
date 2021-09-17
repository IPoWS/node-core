package nodes

import (
	"os"
	"sync"

	"github.com/sirupsen/logrus"
)

type Nodes struct {
	NodesProto
	FileMu sync.RWMutex
	MemMu  sync.RWMutex
}

func (m *Nodes) Save(nodesfile string) error {
	if nodesfile == "" {
		nodesfile = "./nodes"
	}
	m.MemMu.RLock()
	data, err := m.Marshal()
	m.MemMu.RUnlock()
	if err == nil {
		m.FileMu.Lock()
		err = os.WriteFile(nodesfile, data, 0644)
		m.FileMu.Unlock()
	}
	return err
}

func (m *Nodes) Load(nodesfile string) error {
	if nodesfile == "" {
		nodesfile = "./nodes"
	}
	m.FileMu.RLock()
	data, err := os.ReadFile(nodesfile)
	m.FileMu.RUnlock()
	logrus.Debugf("[nodes] load %d bytes from file.", len(data))
	if err == nil && data != nil && len(data) > 0 {
		m.MemMu.Lock()
		err = m.Unmarshal(data)
		m.MemMu.Unlock()
	} else {
		m.Clear()
	}
	return err
}

func (m *Nodes) Clear() {
	m.Hosts = make(map[string]uint64)
	m.Ip64S = make(map[uint64]string)
	m.Nodes = make(map[string]string)
	m.Delay = make(map[uint64]uint64)
	m.Names = make(map[uint64]string)
	logrus.Debugln("[nodes] clear node.")
}

func (m *Nodes) CopyNodes() map[string]string {
	m.MemMu.RLock()
	ret := make(map[string]string, len(m.Nodes))
	for k, v := range m.Nodes {
		ret[k] = v
	}
	m.MemMu.RUnlock()
	return ret
}

func (m *Nodes) CopyIp64S() map[uint64]string {
	m.MemMu.RLock()
	ret := make(map[uint64]string, len(m.Ip64S))
	for k, v := range m.Ip64S {
		ret[k] = v
	}
	m.MemMu.RUnlock()
	return ret
}

func (m *Nodes) CopyHosts() map[string]uint64 {
	m.MemMu.RLock()
	ret := make(map[string]uint64, len(m.Hosts))
	for k, v := range m.Hosts {
		ret[k] = v
	}
	m.MemMu.RUnlock()
	return ret
}

func (m *Nodes) CopyDelay() map[uint64]uint64 {
	m.MemMu.RLock()
	ret := make(map[uint64]uint64, len(m.Delay))
	for k, v := range m.Delay {
		ret[k] = v
	}
	m.MemMu.RUnlock()
	return ret
}

func (m *Nodes) AddNode(host string, ent string, ip uint64, name string, delay uint64) {
	logrus.Debugln("[nodes] add node", host, ent, ip, name, delay)
	m.MemMu.Lock()
	m.Nodes[host] = ent
	m.Ip64S[ip] = host
	m.Hosts[host] = ip
	m.Delay[ip] = delay
	m.Names[ip] = name
	m.MemMu.Unlock()
}

func (m *Nodes) DelNodeByIP(ip uint64) {
	m.MemMu.Lock()
	host, ok := m.Ip64S[ip]
	if ok {
		delete(m.Nodes, host)
		delete(m.Hosts, host)
		delete(m.Ip64S, ip)
		delete(m.Delay, ip)
		delete(m.Names, ip)
	}
	m.MemMu.Unlock()
}

func (m *Nodes) IsIp64InNodes(ip uint64) bool {
	m.MemMu.RLock()
	_, ok := m.Ip64S[ip]
	m.MemMu.RUnlock()
	return ok
}

func (m *Nodes) ParseRawNodes(d []byte) error {
	defer m.MemMu.Unlock()
	m.MemMu.Lock()
	m.Clear()
	return m.Unmarshal(d)
}
