package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	sendmap = make(map[uint64]*websocket.Conn)
	sendmu  sync.RWMutex
)

func saveMap(wsip uint64, conn *websocket.Conn) {
	if wsip > 0 {
		sendmu.RLock()
		_, ok := sendmap[wsip]
		sendmu.RUnlock()
		if !ok {
			sendmu.Lock()
			sendmap[wsip] = conn
			sendmu.Unlock()
		}
	}
}

func delMap(wsip uint64) {
	sendmu.RLock()
	conn, ok := sendmap[wsip]
	sendmu.RUnlock()
	if ok {
		sendmu.Lock()
		delete(sendmap, wsip)
		if conn != nil {
			conn.Close()
		}
		sendmu.Unlock()
	}
}

func copyMap() map[uint64]*websocket.Conn {
	sendmu.RLock()
	ret := make(map[uint64]*websocket.Conn, len(sendmap))
	for k, v := range sendmap {
		ret[k] = v
	}
	sendmu.RUnlock()
	return ret
}
