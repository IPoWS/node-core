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
			for i := range NodesList.CopyIp64S() {
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
