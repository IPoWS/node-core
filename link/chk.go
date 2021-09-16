package link

import (
	"time"

	"github.com/sirupsen/logrus"
)

var (
	lastalive = make(map[uint64]uint64)
)

func SetAlive(ip uint64) {
	lastalive[ip] = uint64(time.Now().UnixNano())
}

func isLinkAlive(host string, ent string, ip uint64) bool {
	now := uint64(time.Now().UnixNano())
	t, ok := lastalive[ip]
	if ok {
		return now-t <= 65536*1000000
	} else {
		wsip, _, err := InitLink("ws://"+host+"/"+ent, ip)
		if err != nil || (ip != 0 && wsip != ip) {
			return false
		}
		return true
	}
}

func startCheck() {
	go func() {
		n := NodesList.CopyNodes()
		for ip, host := range NodesList.CopyIp64S() {
			if host != "" && !isLinkAlive(host, n[host], ip) {
				NodesList.DelNodeByIP(ip)
				logrus.Infof("[checklink] del %x -> %s.", ip, host)
			}
		}
		SaveNodesBack()
		t := time.NewTicker(time.Millisecond * 32768)
		for range t.C {
			logrus.Info("[checklink] send hello started.")
			for i := range NodesList.CopyIp64S() {
				err := SendHello(i)
				if err != nil {
					DelConn(i)
				} else {
					SetAlive(i)
				}
			}
			logrus.Info("[checklink] send hello finished.")
		}
	}()
}
