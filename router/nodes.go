package router

import (
	"sync"

	"github.com/IPoWS/node-core/data/nodes"
)

var (
	Allnodes *nodes.Nodes
	Nodesmu  sync.RWMutex
)

func init() {
	Allnodes.Hosts = make(map[string]uint64)
	Allnodes.Ip64S = make(map[uint64]string)
	Allnodes.Nodes = make(map[string]string)
}

func ParseRawNodes(d []byte) error {
	defer Nodesmu.Unlock()
	Nodesmu.Lock()
	Allnodes = new(nodes.Nodes)
	Allnodes.Nodes = make(map[string]string)
	return Allnodes.Unmarshal(d)
}

func AddNode(host string, ent string, ip uint64) {
	Nodesmu.Lock()
	Allnodes.Nodes[host] = ent
	Allnodes.Ip64S[ip] = host
	Allnodes.Hosts[host] = ip
	Nodesmu.Unlock()
}

func DelNode(host string) {
	Nodesmu.Lock()
	_, ok := Allnodes.Nodes[host]
	if ok {
		delete(Allnodes.Nodes, host)
		ip, ok := Allnodes.Hosts[host]
		if ok {
			delete(Allnodes.Hosts, host)
			delete(Allnodes.Ip64S, ip)
		}
	}
	Nodesmu.Unlock()
}

func SaveNodes(nodesfile string) {
	Allnodes.Save(nodesfile)
}

func LoadNodes(nodesfile string) {
	Allnodes.Load(nodesfile)
}
