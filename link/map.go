package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	connmap = make(map[uint64]*WSNode)
	connmu  sync.RWMutex
)

func saveMap(wsip uint64, conn *websocket.Conn, mt int) {
	connmu.Lock()
	connmap[wsip] = new(WSNode)
	connmap[wsip].conn = conn
	connmu.Unlock()
}

func updateMt(wsip uint64, mt int) {
	if wsip > 0 {
		connmu.Lock()
		connmap[wsip].mt = mt
		connmu.Unlock()
	}
}
