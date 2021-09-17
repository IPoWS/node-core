package router

import (
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

var (
	table = new(transTable)
)

func init() {
	table.init()
}

func AddItem(to uint64, next uint64, delay100us uint16, conn *websocket.Conn) {
	logrus.Infof("[router] add %x, next hop %x, delay %d * 100us.", to, next, delay100us)
	table.add(&TransItem{to, next, delay100us, conn})
}

func DelItem(to uint64) {
	logrus.Infof("[router] del %x.", to)
	table.del(to)
}

func NextHop(to uint64) *TransItem {
	return table.nextHop(to)
}

func NearMe() []*TransItem {
	return table.near()
}

func AllNeighbors() []*TransItem {
	return table.all()
}

func IsIn(ip uint64) bool {
	return table.isIn(ip)
}
