package link

import (
	"github.com/sirupsen/logrus"
)

func isLinkAlive(host string, ent string, ip uint64) bool {
	wsip, _, err := InitLink("ws://"+host+"/"+ent, ip)
	if err != nil || (ip != 0 && wsip != ip) {
		return false
	}
	return true
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
	}()
}
