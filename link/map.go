package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	connmap = make(map[uint64]*websocket.Conn)
	connmu  sync.RWMutex
)

func saveMap(wsip uint64, conn *websocket.Conn) {
	if wsip > 0 {
		connmu.RLock()
		_, ok := connmap[wsip]
		connmu.RUnlock()
		if !ok {
			connmu.Lock()
			connmap[wsip] = conn
			connmu.Unlock()
		}
	}
}

func delMap(wsip uint64) {
	conn, ok := connmap[wsip]
	if ok {
		connmu.Lock()
		delete(connmap, wsip)
		if conn != nil {
			conn.Close()
		}
		connmu.Unlock()
	}
}
