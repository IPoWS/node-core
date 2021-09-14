package link

import (
	"sync"
	"time"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/sirupsen/logrus"
)

var (
	NodesList *nodes.Nodes
	nfile     string
	newnodes  *nodes.Nodes
	nnt       = time.NewTicker(time.Minute)
	nnmu      sync.Mutex
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

func sendNewNodes() {
	newnodes.MemMu.Lock()
	data, err := newnodes.Marshal()
	newnodes.Clear()
	newnodes.MemMu.Unlock()
	if err == nil {
		for i := range NodesList.CopyIp64S() {
			Send(i, &data, ip64.NodesType)
			logrus.Infof("[link] send new node info %x.", i)
		}
	}
}

// RegisterNode 注册新的节点
func RegisterNode(host string, ent string, ip uint64, name string, delay uint64) {
	newnodes.AddNode(host, ent, ip, name, delay)
}

func startDeliverNewNodes() {
	go func() {
		for {
			select {
			case <-nnt.C:
				if len(newnodes.Names) > 0 {
					sendNewNodes()
				}
			default:
				if len(newnodes.Names) > 4 {
					sendNewNodes()
				}
			}
		}
	}()
}
