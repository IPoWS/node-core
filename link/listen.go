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

// listen 监听其他节点发来的包
func listen(conn *websocket.Conn) {
	var err error
	for err == nil {
		_, p, err := conn.ReadMessage()
		if err == nil {
			var ip ip64.Ip64
			err := ip.Unmarshal(p)
			if err == nil {
				if ip.To == Mywsip {
					switch ip.Prototype {
					case ip64.HelloType:
						t := time.Now().UnixNano()
						var h hello.Hello
						err = h.Unmarshal(ip.Data)
						delay := t - ip.Time
						logrus.Infof("[listen] recv hello from: %x, to: %x, delay: %d ns.", ip.From, ip.To, delay)
						if err == nil && delay > 0 && ip.From > 0 && ip.To > 0 {
							saveMap(ip.From, conn)
							router.AddItem(ip.From, ip.From, uint16(delay/10000))
							SetAlive(ip.From)
						}
						if err == nil {
							if Mywsip == 0 {
								Mywsip = ip.To
								myhello.Mask = h.Mask
								logrus.Infof("[listen] set my ip: %x with mask %x.", Mywsip, h.Mask)
								saveMap(Mywsip, conn)
							}
							if Mywsip > 0 {
								SendHello(Mywsip)
							}
						}
					case ip64.NodesType: // 在地址列表更新后
						logrus.Info("[listen] recv nodes.")
						var newnodes nodes.Nodes
						newnodes.Unmarshal(ip.Data)
						for wsip, host := range newnodes.Ip64S {
							ent := newnodes.Nodes[host]
							if isLinkAlive(host, ent, wsip) {
								NodesList.AddNode(host, ent, wsip, newnodes.Names[wsip], uint64((ip.Time-time.Now().UnixNano())/10000))
							} else {
								router.AddItem(wsip, ip.From, uint16((int64(newnodes.Delay[wsip])+ip.Time-time.Now().UnixNano())/10000))
							}
						}
					}
				} else {
					logrus.Info("[listen] forward pack from %x to %x.", ip.From, ip.To)
					Forward(router.NextHop(ip.To), &ip)
				}
			}
		} else {
			logrus.Errorf("[listen] %v", err)
			err = nil
		}
	}
	logrus.Errorf("[listen] %v", err)
	conn.Close()
}

// SendHello 从自身发送 hello 给对方
func SendHello(to uint64) error {
	h := myhello
	return sendHello(to, &h)
}

func DelConn(wsip uint64) {
	logrus.Infof("[delconn] %x is unreachable and del it from table.", wsip)
	router.DelItem(wsip)
	NodesList.DelNodeByIP(wsip)
	delMap(wsip)
	SaveNodesBack()
}

// sendHello 发送 hello 给对方
func sendHello(wsip uint64, h *hello.Hello) error {
	connmu.RLock()
	wsn, ok := connmap[wsip]
	connmu.RUnlock()
	if ok {
		data, err := h.Marshal()
		if err == nil {
			var ip ip64.Ip64
			ip.Pack(Mywsip, wsip, &data, ip64.HelloType)
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
		ip.Pack(Mywsip, adviceip, &data, ip64.HelloType)
		logrus.Info("[sendHello] send hello.")
		err = ip.Send(conn, websocket.BinaryMessage)
	}
	if err != nil {
		logrus.Errorf("[sendHello] %v", err)
	}
	return err
}
