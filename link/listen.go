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
			err = ip.Unmarshal(p)
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
							logrus.Infof("[listen.hello] recv from: %x, to: %x, delay: %d ns.", ip.From, ip.To, delay)
							if err == nil {
								if ip.From > 0 && ip.To > 0 && Mywsip > 0 {
									// saveMap(ip.From, conn)
									router.AddItem(ip.From, ip.From, uint16(delay/100000))
									NodesList.AddNode(h.Host, h.Entry, ip.From, h.Name, uint64(delay))
									registerNode(ip.From)
								} else {
									if Mywsip == 0 {
										Mywsip = ip.To
										myhello.Mask = h.Mask
										logrus.Infof("[listen.hello] set my ip: %x with mask %x.", Mywsip, h.Mask)
										// saveMap(Mywsip, conn)
										router.AddItem(ip.To, ip.To, uint16(delay/100000))
										NodesList.AddNode(h.Host, h.Entry, ip.To, h.Name, uint64(delay))
										registerNode(ip.To)
									}
								}
								if h.Isinit {
									logrus.Infoln("[listen.hello] recv init.")
									if ip.From == 0 {
										sendmu.Lock()
										sendmap[Mywsip] = sendmap[0]
										sendmu.Unlock()
										SendHello(Mywsip) // send to me
										InitLink("ws://"+conn.RemoteAddr().String()+"/"+h.Entry, Mywsip)
									} else {
										if !isInMap(ip.From) {
											InitLink("ws://"+conn.RemoteAddr().String()+"/"+h.Entry, ip.From)
										}
										SendHello(ip.From)
									}
								}
							} else {
								logrus.Errorln("[listen.hello] unmashal err: ", err)
								err = nil
							}
						case ip64.NodesType: // 在地址列表更新后
							logrus.Info("[listen.nodes] recv nodes.")
							var newnodes nodes.Nodes
							newnodes.Unmarshal(ip.Data)
							for wsip, host := range newnodes.Ip64S {
								if wsip == Mywsip {
									logrus.Infoln("[listen.nodes] check own node ip.")
									if myhello.Host == "" {
										myhello.Host = host
										listener, err := net.Listen("tcp", host)
										if err == nil {
											go logrus.Fatal(http.Serve(listener, nil))
											logrus.Infof("[listen.nodes] start listening %s.", host)
										} else {
											logrus.Errorf("[listen.nodes] listen %s err: %v.", host, err)
										}
									}
								} else {
									ent := newnodes.Nodes[host]
									alive := isLinkAlive(host, ent, wsip)
									if alive {
										NodesList.AddNode(host, ent, wsip, newnodes.Names[wsip], uint64(delay))
										logrus.Infof("[listen.nodes] add node %x directly.", wsip)
									}
									relay := int64(newnodes.Delay[wsip]) + delay
									if (alive && relay < int64(NodesList.Delay[wsip])) || (!alive && relay < int64(time.Second)) {
										NodesList.AddNode(host, ent, wsip, newnodes.Names[wsip], uint64(relay))
										router.AddItem(wsip, ip.From, uint16(relay/100000))
										logrus.Infof("[listen.nodes] add node %x through %x, delay %d ms.", wsip, ip.From, relay/10)
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
		}
	}
	logrus.Errorf("[listen] %v", err)
	conn.Close()
}
