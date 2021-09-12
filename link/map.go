package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	connmap = make(map[uint64]*WSNode)
	mapmu   sync.RWMutex
)

func saveMap(wsip uint64, conn *websocket.Conn, mt int, delay int64) {
	mapmu.Lock()
	connmap[wsip] = new(WSNode)
	connmap[wsip].conn = conn
	connmap[wsip].mt = mt
	connmap[wsip].delay = delay
	mapmu.Unlock()
}
