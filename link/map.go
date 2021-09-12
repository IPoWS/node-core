package link

import (
	"sync"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/gorilla/websocket"
)

var (
	connmap = make(map[uint32]*WSNode)
	mapmu   sync.RWMutex
)

func saveMap(h *hello.Hello, conn *websocket.Conn, mt int, delay int64) {
	mapmu.Lock()
	connmap[h.Wsnetaddr] = new(WSNode)
	connmap[h.Wsnetaddr].conn = conn
	connmap[h.Wsnetaddr].mt = mt
	connmap[h.Wsnetaddr].delay = delay
	mapmu.Unlock()
}
