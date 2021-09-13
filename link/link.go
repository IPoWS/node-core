package link

import (
	"fmt"
	"time"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func Send(to uint64, fromport uint16, toport uint16, data *[]byte, prototype uint32) error {
	connmu.RLock()
	wsn, ok := connmap[to]
	connmu.RUnlock()
	if ok {
		var ip ip64.Ip64
		ip.Pack(mywsip, to, fromport, toport, data, prototype)
		logrus.Info("[Send] link send %d bytes to %x:%d from %d.", len(*data), to, toport, fromport)
		return ip.Send(wsn, websocket.BinaryMessage)
	}
	return fmt.Errorf("dest %x unreachable.", to)
}

func Forward(to uint64, ip *ip64.Ip64) error {
	connmu.RLock()
	wsn, ok := connmap[to]
	connmu.RUnlock()
	if ok {
		return ip.Send(wsn, websocket.BinaryMessage)
	}
	return fmt.Errorf("dest %x unreachable.", to)
}

func StartCheck(m *nodes.Nodes) {
	go func() {
		n := m.CopyNodes()
		for ip, host := range m.CopyIp64S() {
			wsip, _, err := InitLink("ws://"+host+"/"+n[host], ip)
			if err != nil || wsip != ip {
				router.DelNodeByIP(ip)
				logrus.Infof("[linkcheck] del %x -> %s.", ip, host)
			}
		}
		router.SaveNodesBack()
		t := time.NewTicker(time.Millisecond * 327)
		t2 := time.NewTicker(time.Millisecond * 655)
		select {
		case <-t.C:
			logrus.Info("[checklick] send hello.")
			for i := range m.CopyIp64S() {
				SendHello(i)
			}
		case <-t2.C:
			logrus.Info("[checklick] check alive.")
			now := uint64(time.Now().UnixNano())
			for i, t := range m.CopyTimes() {
				if now-t > 65536*1000000 {
					DelConn(i)
				}
			}
		}
	}()
}
