package link

import (
	"github.com/IPoWS/node-core/router"
	"github.com/sirupsen/logrus"
)

func isLinkAlive(host string, ent string, ip uint64) bool {
	if router.IsIn(ip) {
		return true
	}
	wsip, _, err := InitLink("ws://"+host+"/"+ent, ip, true)
	if err != nil || (ip != 0 && wsip != ip) {
		logrus.Debugf("[link.chk] link to host %s with ip %x (advice %x) err: %v.", host, wsip, ip, err)
		return false
	}
	logrus.Debugf("[link.chk] link to host %s with ip %x (advice %x) alive.")
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
