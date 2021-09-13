package router

import (
	"sync"

	"github.com/IPoWS/node-core/data/nodes"
)

var (
	Allnodes *nodes.Nodes
	Nodesmu  sync.RWMutex
	nfile    string
)

func init() {
	Allnodes = new(nodes.Nodes)
}

func ParseRawNodes(d []byte) error {
	defer Nodesmu.Unlock()
	Nodesmu.Lock()
	return Allnodes.Unmarshal(d)
}

func AddNode(host string, ent string, ip uint64, time uint64) {
	Nodesmu.Lock()
	Allnodes.Nodes[host] = ent
	Allnodes.Ip64S[ip] = host
	Allnodes.Hosts[host] = ip
	Allnodes.Times[ip] = time
	Nodesmu.Unlock()
}

func FlushAlive(ip uint64, time uint64) {
	Nodesmu.Lock()
	Allnodes.Times[ip] = time
	Nodesmu.Unlock()
}

func DelNodeByHost(host string) {
	Nodesmu.Lock()
	_, ok := Allnodes.Nodes[host]
	if ok {
		delete(Allnodes.Nodes, host)
		ip, ok := Allnodes.Hosts[host]
		if ok {
			delete(Allnodes.Hosts, host)
			delete(Allnodes.Ip64S, ip)
			delete(Allnodes.Times, ip)
		}
	}
	Nodesmu.Unlock()
}

func DelNodeByIP(ip uint64) {
	Nodesmu.Lock()
	host, ok := Allnodes.Ip64S[ip]
	if ok {
		delete(Allnodes.Nodes, host)
		delete(Allnodes.Hosts, host)
		delete(Allnodes.Ip64S, ip)
		delete(Allnodes.Times, ip)
	}
	Nodesmu.Unlock()
}

func SaveNodes(nodesfile string) error {
	return Allnodes.Save(nodesfile)
}

func SaveNodesBack() error {
	return Allnodes.Save(nfile)
}

func LoadNodes(nodesfile string) error {
	nfile = nodesfile
	return Allnodes.Load(nodesfile)
}
