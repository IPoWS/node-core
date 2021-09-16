package link

import (
	"github.com/IPoWS/node-core/router"
	"github.com/sirupsen/logrus"
)

func DelConn(wsip uint64) {
	logrus.Infof("[delconn] %x is unreachable and del it from table.", wsip)
	router.DelItem(wsip)
	NodesList.DelNodeByIP(wsip)
	delMap(wsip)
	SaveNodesBack()
}
