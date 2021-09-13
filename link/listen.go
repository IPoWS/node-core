package link

import (
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	myhello hello.Hello
)

// listen 监听其他节点发来的包
func listen(conn *websocket.Conn) {
	var err error
	for err == nil {
		_, p, err := conn.ReadMessage()
		if err == nil {
			var ip ip64.Ip64
			err := ip.Unmarshal(p)
			if err == nil {
				switch ip.Prototype {
				case ip64.HelloType:
					logrus.Info("[listen] recv hello.")
					t := time.Now().UnixNano()
					var h hello.Hello
					err = h.Unmarshal(ip.Data)
					delay := t - h.Time
					logrus.Infof("[listen] from: %x, to: %x, delay: %d ms.", ip.From, ip.To, delay/1000000)
					if err == nil && delay > 0 && ip.From > 0 && ip.To > 0 {
						saveMap(ip.From, conn)
						router.AddItem(ip.From, ip.From, uint16(delay/1000000))
						h = myhello
						h.Time = time.Now().UnixNano()
						err = sendHello(ip.From, &h)
						if err == nil {
							if mywsip == 0 {
								mywsip = ip.To
								mymask = h.Mask
								logrus.Infof("[listen] set my ip: %x with mask %x.", mywsip, mymask)
							}
							return
						}
					}
				case ip64.NodesType: // 在地址列表更新后
					logrus.Info("[listen] recv nodes.")
					var newnodes nodes.Nodes
					newnodes.Unmarshal(ip.Data)
					for h, e := range newnodes.Nodes {
						if e == "" {
							router.DelNode(h)
						} else {
							router.AddNode(h, e)
							InitLink(h+e, 0)
						}
					}
				}
			}
		}
	}
	logrus.Errorf("[listen] %v", err)
	conn.Close()
}

// sendHello 发送 hello 给对方
func sendHello(wsip uint64, h *hello.Hello) error {
	wsn, ok := connmap[wsip]
	if ok {
		data, err := h.Marshal()
		if err == nil {
			var ip ip64.Ip64
			ip.Pack(mywsip, wsip, 0, 0, &data, uintptr(len(data)), ip64.HelloType)
			logrus.Info("[sendHello] send hello.")
			err = ip.Send(wsn, websocket.BinaryMessage)
		}
		if err != nil {
			logrus.Errorf("[sendHello] %v", err)
		}
		return err
	} else {
		logrus.Infof("[sendHello] destination %x is unreachable.", wsip)
		return nil
	}
}

// sendHelloUnknown 发送 hello 给未知 ip 方
func sendHelloUnknown(conn *websocket.Conn, h *hello.Hello, adviceip uint64) error {
	data, err := h.Marshal()
	if err == nil {
		var ip ip64.Ip64
		ip.Pack(mywsip, adviceip, 0, 0, &data, uintptr(len(data)), ip64.HelloType)
		logrus.Info("[sendHello] send hello.")
		err = ip.Send(conn, websocket.BinaryMessage)
	}
	if err != nil {
		logrus.Errorf("[sendHello] %v", err)
	}
	return err
}
