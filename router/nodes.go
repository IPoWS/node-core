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
	Allnodes = new(nodes.Nodes)
}

func ParseRawNodes(d []byte) error {
	defer Nodesmu.Unlock()
	Nodesmu.Lock()
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

func SaveNodes(nodesfile string) error {
	return Allnodes.Save(nodesfile)
}

func LoadNodes(nodesfile string) error {
	return Allnodes.Load(nodesfile)
}
