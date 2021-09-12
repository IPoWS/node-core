package link

import (
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/gorilla/websocket"
)

type WSNode struct {
	conn  *websocket.Conn
	mt    int
	delay int64
}

var (
	myhello hello.Hello
)

// listenHello 监听其他节点的 hello
func listenHello(conn *websocket.Conn) {
	mt, p, err := conn.ReadMessage()
	t := time.Now().UnixNano()
	if err == nil {
		var h hello.Hello
		err = h.Unmarshal(p)
		delay := t - h.Time
		if err == nil && delay > 0 {
			saveMap(&h, conn, mt, delay)
			h = myhello
			h.Time = time.Now().UnixNano()
			err = sendHello(conn, mt, &h)
			if err == nil {
				return
			}
		}
	}
	conn.Close()
}

// sendHello 发送 hello 给对方
func sendHello(conn *websocket.Conn, mt int, h *hello.Hello) error {
	data, err := h.Marshal()
	if err == nil {
		err = conn.WriteMessage(mt, data)
	}
	return err
}
