package link

import (
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/ip64"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

type WSNode struct {
	conn  *websocket.Conn
	mt    int
	delay int64
}

var (
	myhello hello.Hello
)

// listen 监听其他节点发来的包
func listen(conn *websocket.Conn) {
	var err error
	for err == nil {
		mt, p, err := conn.ReadMessage()
		if err == nil {
			var ip ip64.Ip64
			err := ip.Unmarshal(p)
			if err == nil {
				switch ip.Prototype {
				case ip64.HelloType:
					logrus.Info("[listenHello] recv hello.")
					t := time.Now().UnixNano()
					var h hello.Hello
					err = h.Unmarshal(ip.Data)
					delay := t - h.Time
					if err == nil && delay > 0 {
						saveMap(ip.From, conn, mt, delay)
						h = myhello
						h.Time = time.Now().UnixNano()
						err = sendHello(ip.From, &h)
						if err == nil {
							return
						}
					}
				case ip64.NodesType:
				}
			}
		}
	}
	logrus.Errorf("[listenHello] %v", err)
	conn.Close()
}

// sendHello 发送 hello 给对方
func sendHello(wsip uint64, h *hello.Hello) error {
	wsn, ok := connmap[wsip]
	if ok {
		data, err := h.Marshal()
		if err == nil {
			var ip ip64.Ip64
			ip.Pack(mywsip, wsip, 0, 0, data, uintptr(len(data)), ip64.HelloType)
			logrus.Info("[sendHello] send hello.")
			err = ip.Send(wsn.conn, wsn.mt)
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
func sendHelloUnknown(conn *websocket.Conn, mt int, h *hello.Hello) error {
	data, err := h.Marshal()
	if err == nil {
		var ip ip64.Ip64
		ip.Pack(mywsip, 0, 0, 0, data, uintptr(len(data)), ip64.HelloType)
		logrus.Info("[sendHello] send hello.")
		err = ip.Send(conn, mt)
	}
	if err != nil {
		logrus.Errorf("[sendHello] %v", err)
	}
	return err
}
