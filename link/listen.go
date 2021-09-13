package link

import (
	"sync"
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	myhello  hello.Hello
	alivemap = make(map[uint64]bool)
	alimu    sync.RWMutex
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
					logrus.Infof("[listen] recv hello with %d bytes data.", len(ip.Data))
					t := time.Now().UnixNano()
					var h hello.Hello
					err = h.Unmarshal(ip.Data)
					delay := t - h.Time
					logrus.Infof("[listen] from: %x, to: %x, delay: %d ns.", ip.From, ip.To, delay)
					alimu.Lock()
					alivemap[ip.From] = true
					alimu.Unlock()
					if err == nil && delay > 0 && ip.From > 0 && ip.To > 0 {
						saveMap(ip.From, conn)
						router.AddItem(ip.From, ip.From, uint16(delay/10000))
					}
					if err == nil {
						if mywsip == 0 {
							mywsip = ip.To
							mymask = h.Mask
							logrus.Infof("[listen] set my ip: %x with mask %x.", mywsip, mymask)
							saveMap(mywsip, conn)
						}
						if mywsip > 0 {
							h = myhello
							h.Time = time.Now().UnixNano()
							err = sendHello(mywsip, &h)
						}
					}
				case ip64.NodesType: // 在地址列表更新后
					logrus.Info("[listen] recv nodes.")
					var newnodes nodes.Nodes
					newnodes.Unmarshal(ip.Data)
					for h, e := range newnodes.Nodes {
						if e == "" {
							router.DelNodeByHost(h)
						} else {
							router.AddNode(h, e, newnodes.Hosts[h])
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

// SendHello 从自身发送 hello 给对方
func SendHello(to uint64) error {
	h := myhello
	h.Time = time.Now().UnixNano()
	return sendHello(to, &h)
}

// sendHello 发送 hello 给对方
func sendHello(wsip uint64, h *hello.Hello) error {
	connmu.RLock()
	wsn, ok := connmap[wsip]
	connmu.RUnlock()
	alimu.Lock()
	alivemap[wsip] = false
	alimu.Unlock()
	if ok {
		data, err := h.Marshal()
		if err == nil {
			var ip ip64.Ip64
			ip.Pack(mywsip, wsip, 0, 0, &data, ip64.HelloType)
			logrus.Info("[sendHello] send hello.")
			err = ip.Send(wsn, websocket.BinaryMessage)
		}
		if err != nil {
			logrus.Errorf("[sendHello] %v", err)
		} else {
			go func() {
				// sleep 65.536 s
				time.Sleep(time.Millisecond * 65536)
				alimu.RLock()
				ok = alivemap[wsip]
				alimu.RUnlock()
				if !ok {
					logrus.Infof("[sendHello] %x is unreachable and del it from table.", wsip)
					router.DelItem(wsip)
					router.DelNodeByIP(wsip)
					delMap(wsip)
					router.SaveNodesBack()
				}
			}()
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
		ip.Pack(mywsip, adviceip, 0, 0, &data, ip64.HelloType)
		logrus.Info("[sendHello] send hello.")
		err = ip.Send(conn, websocket.BinaryMessage)
	}
	if err != nil {
		logrus.Errorf("[sendHello] %v", err)
	}
	return err
}
