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
func registerNode(host string, ent string, to uint64, name string, delay uint64) {
	newnodes.AddNode(host, ent, to, name, delay)
	logrus.Infof("[registerNode] %x.", to)
}

func startDeliverNewNodes() {
	go func() {
		for range nnt.C {
			n := router.NearMe()
			if len(n) > 0 {
				for _, i := range n {
					registerNode(myhost, myhello.Entry, i.To, NodesList.Names[i.To], uint64(i.Delay100us)*100000)
				}
				SendNewNodes(newnodes)
			}
		}
	}()
}
