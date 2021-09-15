package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	sendmap = make(map[uint64]*websocket.Conn)
	connmu  sync.RWMutex
)

func saveMap(wsip uint64, conn *websocket.Conn) {
	if wsip > 0 {
		connmu.RLock()
		_, ok := sendmap[wsip]
		connmu.RUnlock()
		if !ok {
			connmu.Lock()
			sendmap[wsip] = conn
			connmu.Unlock()
		}
	}
}

func delMap(wsip uint64) {
	conn, ok := sendmap[wsip]
	if ok {
		connmu.Lock()
		delete(sendmap, wsip)
		if conn != nil {
			conn.Close()
		}
		connmu.Unlock()
	}
}
