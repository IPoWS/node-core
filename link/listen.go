package link

import (
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
	logrus.Infof("[listen] wait for msg.")
	_, p, err := conn.ReadMessage()
	logrus.Infof("[listen] msg arrived.")
	if err == nil {
		var ip ip64.Ip64
		err = ip.Unmarshal(p)
		sendhellobackto := ip.From
		if err == nil {
			logrus.Infof("[listen] recv from: %x, to: %x.", ip.From, ip.To)
			if Mywsip == 0 || (Mywsip != 0 && ip.To == Mywsip) {
				t := time.Now().UnixNano()
				delay := t - ip.Time
				if delay < int64(time.Second*6) && delay > 0 {
					switch uint16(ip.Destproto & 0x0000_ffff) {
					case ip64.HelloType:
						var h hello.Hello
						err = h.Unmarshal(ip.Data)
						logrus.Infof("[listen.hello] delay: %d ns.", delay)
						if err == nil {
							if h.Isinit {
								logrus.Infoln("[listen.hello] recv init.")
								if Mywsip == 0 {
									Mywsip = ip.To
									myhello.Mask = h.Mask
									logrus.Infof("[listen.hello] set my ip: %x with mask %x.", Mywsip, h.Mask)
									// saveMap(Mywsip, conn)
									router.AddItem(ip.To, ip.To, uint16(delay/100000))
									NodesList.AddNode(conn.RemoteAddr().String(), h.Entry, ip.To, h.Name, uint64(delay))
									registerNode(ip.To)
									if ip.From == 0 { // 自分配ip，是nps
										sendmu.Lock()
										sendmap[Mywsip] = conn
										sendmu.Unlock()
										sendhellobackto = Mywsip
									}
								}
								if ip.From != 0 { // 是其它node建立的链接，建立一条反向链接以send
									sendmu.Lock()
									_, ok := sendmap[ip.From]
									sendmap[ip.From] = conn
									sendmu.Unlock()
									if !ok {
										InitLink("ws://"+conn.RemoteAddr().String()+"/"+h.Entry, ip.From)
									}
								}
							} else if ip.From > 0 && ip.To > 0 && Mywsip > 0 {
								// saveMap(ip.From, conn)
								router.AddItem(ip.From, ip.From, uint16(delay/100000))
								NodesList.AddNode(conn.RemoteAddr().String(), h.Entry, ip.From, h.Name, uint64(delay))
								registerNode(ip.From)
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
							if wsip != Mywsip {
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
						destport := uint16((ip.Destproto & 0xffff_0000) >> 16)
						srcport := uint16((uint32(ip.Srcttl) & 0xffff_0000) >> 16)
						logrus.Infof("[listen.data] recv data from port %d to %d.", srcport, destport)
						upper.Recv(srcport, destport, &ip.Data)
					}
				} else {
					logrus.Infof("[listen] delay of package from %x is invalid.", ip.From)
				}
			} else {
				logrus.Infof("[listen] forward pack from %x to %x.", ip.From, ip.To)
				Forward(router.NextHop(ip.To), &ip)
			}
			SendHello(sendhellobackto, conn)
		}
	}
}
