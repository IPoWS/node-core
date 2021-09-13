package router

import (
	"sync"

	"github.com/IPoWS/node-core/data/nodes"
)

var (
	allnodes *nodes.Nodes
	nodesmu  sync.RWMutex
)

func ParseRawNodes(d []byte) error {
	defer nodesmu.Unlock()
	nodesmu.Lock()
	allnodes = new(nodes.Nodes)
	allnodes.Nodes = make(map[string]string)
	return allnodes.Unmarshal(d)
}

func AddNode(host string, ent string) {
	nodesmu.Lock()
	allnodes.Nodes[host] = ent
	nodesmu.Unlock()
}

func DelNode(host string) {
	nodesmu.Lock()
	_, ok := allnodes.Nodes[host]
	if ok {
		delete(allnodes.Nodes, host)
	}
	nodesmu.Unlock()
}
