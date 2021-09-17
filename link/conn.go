package link

import (
	"github.com/IPoWS/node-core/router"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func DelConn(wsip uint64) {
	logrus.Infof("[delconn] %x is unreachable and del it from table.", wsip)
	router.DelItem(wsip)
	NodesList.DelNodeByIP(wsip)
	SaveNodesBack()
}

func AddDirectConn(to uint64, host string, ent string, name string, delay uint64, mask uint64, conn *websocket.Conn) {
	router.AddItem(to&mask, to&mask, uint16(delay/100000), conn)
	NodesList.AddNode(host, ent, to, name, delay)
	registerNode(host, ent, to, name, delay)
}
