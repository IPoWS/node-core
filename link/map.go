package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	connmap = make(map[uint64]*WSNode)
	connmu  sync.RWMutex
)

func checkNotNil(wsip uint64) {
	_, ok := connmap[wsip]
	if !ok {
		connmap[wsip] = new(WSNode)
	}
}

func saveMap(wsip uint64, conn *websocket.Conn, mt int) {
	if wsip > 0 {
		connmu.Lock()
		checkNotNil(wsip)
		connmap[wsip].conn = conn
		connmu.Unlock()
	}
}

func updateMt(wsip uint64, mt int) {
	if wsip > 0 {
		connmu.Lock()
		checkNotNil(wsip)
		connmap[wsip].mt = mt
		connmu.Unlock()
	}
}
