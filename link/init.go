package link

import (
	"net/http"
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	npsurl string
)

// SetNPSUrl 设置NPS服务器地址
func SetNPSUrl(url string) {
	npsurl = url
}

// InitLink 初始化连接 返回 conn, messageType, delay, error
func InitLink(url string) (conn *websocket.Conn, mt int, delay int64, err error) {
	log.Printf("[initlink] connecting to %s", url)
	conn, _, err = websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return
	}
	t := time.Now().UnixNano()
	myhello.Time = t
	sendHello(conn, 0, &myhello)
	mt, p, err := conn.ReadMessage()
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return
	}
	var hello hello.Hello
	err = hello.Unmarshal(p)
	if err != nil {
		log.Errorf("[initlink] parse hello err: %v", err)
		return
	}
	delay = hello.Time - t
	if delay <= 0 {
		log.Errorf("[initlink] tr: %v, t: %v", hello.Time, t)
		return
	}
	saveMap(&hello, conn, mt, delay)
	log.Printf("[initlink] %s 链接测试成功，延时%vns", url, delay)
	return
}

var upgrader = websocket.Upgrader{}

func InitEntry(ent string) {
	myhello = hello.Hello{
		Entry: ent,
	}
	http.HandleFunc("/"+ent, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err == nil {
			go listenHello(conn)
		}
	})
}
