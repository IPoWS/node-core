package link

import (
	"time"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/sirupsen/logrus"
)

var (
	lastalive = make(map[uint64]uint64)
)

func SetAlive(ip uint64) {
	lastalive[ip] = uint64(time.Now().UnixNano())
}

func isLinkAlive(host string, ent string, ip uint64) bool {
	wsip, _, err := InitLink("ws://"+host+"/"+ent, ip)
	if err != nil || (ip != 0 && wsip != ip) {
		return false
	}
	return true
}

func startCheck(m *nodes.Nodes) {
	go func() {
		n := m.CopyNodes()
		for ip, host := range m.CopyIp64S() {
			if !isLinkAlive(host, n[host], ip) {
				NodesList.DelNodeByIP(ip)
				logrus.Infof("[checklink] del %x -> %s.", ip, host)
			}
		}
		SaveNodesBack()
		t := time.NewTicker(time.Millisecond * 32768)
		for range t.C {
			for i := range m.CopyIp64S() {
				SendHello(i)
				time.Sleep(time.Millisecond * 10)
			}
			logrus.Info("[checklink] send hello finished.")
			time.Sleep(time.Millisecond * 8192)
			now := uint64(time.Now().UnixNano())
			for i, t := range lastalive {
				if now-t > 65536*1000000 {
					DelConn(i)
				}
			}
			logrus.Info("[checklink] check alive finished.")
		}
	}()
}
