package link

import (
	"io"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func Register() error {
	q := npsurl + "?ent=" + myhello.Entry + "&name=" + myhello.Name
	conn, resp, err := websocket.DefaultDialer.Dial(q, nil)
	logrus.Info("[link.Register] register to ", q)
	if err == nil {
		go listen(conn)
		data, err := io.ReadAll(resp.Body)
		if err == nil {
			NodesList.ParseRawNodes(data)
			startCheck()
		} else {
			logrus.Errorln("[link.Register] read body: ", err)
		}
	} else {
		logrus.Errorln("[link.Register] dial: ", err)
	}
	return err
}

// add: full host+entry; del: host+null entry
func NotifyChange(n *nodes.Nodes) {
	data, err := n.Marshal()
	if err == nil {
		for to, wsn := range copyMap() {
			if to > 0 && wsn != nil {
				var ip ip64.Ip64
				ip.Pack(Mywsip, to, &data, ip64.NodesType, 0, 0)
				ip.Send(wsn, websocket.BinaryMessage)
			}
		}
	}
}
