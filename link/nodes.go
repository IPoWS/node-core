package link

import (
	"time"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/sirupsen/logrus"
)

var (
	NodesList *nodes.Nodes
	nfile     string
	newnodes  *nodes.Nodes
	nnt       = time.NewTicker(time.Minute)
)

func init() {
	NodesList = new(nodes.Nodes)
	newnodes = new(nodes.Nodes)
	newnodes.Clear()
	startDeliverNewNodes()
}

func SaveNodes(nodesfile string) error {
	return NodesList.Save(nodesfile)
}

func SaveNodesBack() error {
	return NodesList.Save(nfile)
}

func LoadNodes(nodesfile string) error {
	nfile = nodesfile
	return NodesList.Load(nodesfile)
}

func SendNewNodes(newnodes *nodes.Nodes) {
	newnodes.MemMu.Lock()
	data, err := newnodes.Marshal()
	newnodes.Clear()
	newnodes.MemMu.Unlock()
	if err == nil {
		for i := range NodesList.CopyIp64S() {
			Send(i, &data, ip64.NodesType, 0, 0)
		}
	}
}

// registerNode 注册新的节点到newnodes以便广播
func registerNode(ip uint64, delay uint64) {
	host := NodesList.Ip64S[ip]
	newnodes.AddNode(host, NodesList.Nodes[host], ip, NodesList.Names[ip], delay)
	logrus.Infof("[registerNode] %x.", ip)
}

func startDeliverNewNodes() {
	go func() {
		for range nnt.C {
			n := router.NearMe()
			if n != nil && len(n) > 0 {
				for _, i := range n {
					registerNode(i.To, uint64(i.Delay100us)*100000)
				}
				SendNewNodes(newnodes)
			}
		}
	}()
}
