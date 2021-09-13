package link

import (
	"io"
	"net/http"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/sirupsen/logrus"
)

func Register(ent string) {
	resp, err := http.Get(npsurl + "?ent=" + ent)
	if err == nil {
		data, err := io.ReadAll(resp.Body)
		if err == nil {
			router.ParseRawNodes(data)
		}
	}
	if err != nil {
		logrus.Errorf("[RegisterAndGetNodes] %v", err)
	}
}

// add: full host+entry; del: host+null entry
func NotifyChange(n *nodes.Nodes) {
	data, err := n.Marshal()
	if err == nil {
		for to, wsn := range connmap {
			if to > 0 && wsn.conn != nil {
				var ip ip64.Ip64
				ip.Pack(mywsip, to, 0, 0, &data, uintptr(len(data)), ip64.NodesType)
				ip.Send(wsn.conn, wsn.mt)
			}
		}
	}
}
