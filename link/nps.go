package link

import (
	"io"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func Register(ent string, name string) error {
	conn, resp, err := websocket.DefaultDialer.Dial(npsurl+"?ent="+ent+"&name="+name, nil)
	go listen(conn)
	if err == nil {
		data, err := io.ReadAll(resp.Body)
		if err == nil {
			NodesList.ParseRawNodes(data)
			startCheck(NodesList)
		}
	}
	if err != nil {
		logrus.Errorf("[RegisterAndGetNodes] %v", err)
	}
	return err
}

// add: full host+entry; del: host+null entry
func NotifyChange(n *nodes.Nodes) {
	data, err := n.Marshal()
	if err == nil {
		for to, wsn := range connmap {
			if to > 0 && wsn != nil {
				var ip ip64.Ip64
				ip.Pack(Mywsip, to, &data, ip64.NodesType)
				ip.Send(wsn, websocket.BinaryMessage)
			}
		}
	}
}
