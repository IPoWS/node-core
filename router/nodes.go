package router

import (
	"sync"

	"github.com/IPoWS/node-core/data/nodes"
)

var (
	allnodes *nodes.Nodes
	nodesmu  sync.RWMutex
)

func init() {
	allnodes.Hosts = make(map[string]uint64)
	allnodes.Ip64S = make(map[uint64]string)
	allnodes.Nodes = make(map[string]string)
}

func ParseRawNodes(d []byte) error {
	defer nodesmu.Unlock()
	nodesmu.Lock()
	allnodes = new(nodes.Nodes)
	allnodes.Nodes = make(map[string]string)
	return allnodes.Unmarshal(d)
}

func AddNode(host string, ent string, ip uint64) {
	nodesmu.Lock()
	allnodes.Nodes[host] = ent
	allnodes.Ip64S[ip] = host
	allnodes.Hosts[host] = ip
	nodesmu.Unlock()
}

func DelNode(host string) {
	nodesmu.Lock()
	_, ok := allnodes.Nodes[host]
	if ok {
		delete(allnodes.Nodes, host)
		ip, ok := allnodes.Hosts[host]
		if ok {
			delete(allnodes.Hosts, host)
			delete(allnodes.Ip64S, ip)
		}
	}
	nodesmu.Unlock()
}

func SaveNodes(nodesfile string) {
	allnodes.Save(nodesfile)
}

func LoadNodes(nodesfile string) {
	allnodes.Load(nodesfile)
}
