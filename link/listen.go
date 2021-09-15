package link

import (
	"net"
	"net/http"
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/IPoWS/node-core/upper"
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
				if ip.From > 0 {
					SetAlive(ip.From)
				}
				if Mywsip == 0 || (Mywsip != 0 && ip.To == Mywsip) {
					t := time.Now().UnixNano()
					delay := t - ip.Time
					if delay < int64(time.Second*6) && delay > 0 {
						switch ip.Prototype & 0x0000_ffff {
						case ip64.HelloType:
							var h hello.Hello
							err = h.Unmarshal(ip.Data)
							logrus.Infof("[listen] recv hello from: %x, to: %x, delay: %d ns.", ip.From, ip.To, delay)
							if err == nil && ip.From > 0 && ip.To > 0 {
								saveMap(ip.From, conn)
								router.AddItem(ip.From, ip.From, uint16(delay/100000))
								NodesList.AddNode(h.Host, h.Entry, ip.From, h.Name, uint64(delay))
								registerNode(ip.From)
							}
							if err == nil {
								if Mywsip == 0 {
									Mywsip = ip.To
									myhello.Mask = h.Mask
									logrus.Infof("[listen] set my ip: %x with mask %x.", Mywsip, h.Mask)
									saveMap(Mywsip, conn)
									router.AddItem(ip.To, ip.To, uint16(delay/100000))
									NodesList.AddNode(h.Host, h.Entry, ip.To, h.Name, uint64(delay))
									registerNode(ip.To)
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
								if wsip == Mywsip {
									myhello.Host = host
									listener, err := net.Listen("tcp", host)
									if err == nil {
										go logrus.Fatal(http.Serve(listener, nil))
										logrus.Infof("[listen] start listening %s.", host)
									} else {
										logrus.Infof("[listen] listen %s err: %v.", host, err)
										err = nil
									}
								} else {
									ent := newnodes.Nodes[host]
									alive := isLinkAlive(host, ent, wsip)
									if alive {
										NodesList.AddNode(host, ent, wsip, newnodes.Names[wsip], uint64(delay))
										logrus.Infof("[listen] add node %x directly.", wsip)
									}
									relay := int64(newnodes.Delay[wsip]) + delay
									if (alive && relay < int64(NodesList.Delay[wsip])) || (!alive && relay < int64(time.Second)) {
										NodesList.AddNode(host, ent, wsip, newnodes.Names[wsip], uint64(relay))
										router.AddItem(wsip, ip.From, uint16(relay/100000))
										logrus.Infof("[listen] add node %x through %x, delay %d ms.", wsip, ip.From, relay/10)
									}
								}
							}
						case ip64.DataType:
							port := (ip.Prototype & 0xffff_0000) >> 16
							upper.Recv(uint16(port), &ip.Data)
						}
					} else {
						logrus.Infof("[listen] delay of package from %x is invalid.", ip.From)
					}
				} else {
					logrus.Infof("[listen] forward pack from %x to %x.", ip.From, ip.To)
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
