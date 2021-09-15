package link

import (
	"sync"

	"github.com/gorilla/websocket"
)

var (
	sendmap = make(map[uint64]*websocket.Conn)
	sendmu  sync.RWMutex
)

func isInMap(wsip uint64) bool {
	sendmu.RLock()
	_, ok := sendmap[wsip]
	sendmu.RUnlock()
	return ok
}

func readMap(wsip uint64) (conn *websocket.Conn, ok bool) {
	sendmu.RLock()
	conn, ok = sendmap[wsip]
	sendmu.RUnlock()
	return
}

func saveMap(wsip uint64, conn *websocket.Conn) {
	sendmu.Lock()
	sendmap[wsip] = conn
	sendmu.Unlock()
}

func delMap(wsip uint64) {
	conn, ok := readMap(wsip)
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
