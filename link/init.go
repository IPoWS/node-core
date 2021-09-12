package link

import (
	"strconv"
	"time"

	"github.com/IPoWS/node-core/data"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

// InitLink 初始化连接 返回 conn, messageType, delay, error
func InitLink(url string) (conn *websocket.Conn, mt int, delay int64, err error) {
	log.Printf("[initlink] connecting to %s", url)
	t := time.Now().UnixNano()
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return
	}
	mt, p, err := conn.ReadMessage()
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return
	}
	recvt, err := strconv.ParseInt(data.Bytes2str(p), 10, 64)
	if err != nil {
		log.Errorf("[initlink] parse int err: %v", err)
		return
	}
	delay = recvt - t
	if delay <= 0 {
		log.Errorf("[initlink] tr: %v, t: %v", recvt, t)
		return
	}
	log.Printf("[initlink] %s 链接测试成功，延时%vns", url, delay)
	return
}
