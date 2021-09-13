package link

import (
	"net/http"
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var (
	npsurl string
	mywsip uint64
	mymask uint64
)

// SetNPSUrl 设置NPS服务器地址
func SetNPSUrl(url string) {
	npsurl = url
}

// InitLink 初始化连接 返回 wsip delay error
func InitLink(url string, adviceip uint64) (uint64, int64, error) {
	log.Printf("[initlink] connecting to %s", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return 0, 0, err
	}
	t := time.Now().UnixNano()
	h := myhello
	h.Time = t
	if adviceip > 0 {
		h.Mask = 0xffff_ffff_0000_0000
	}
	sendHelloUnknown(conn, 0, &h, adviceip)
	mt, p, err := conn.ReadMessage()
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return adviceip, 0, err
	}
	var ip ip64.Ip64
	err = ip.Unmarshal(p)
	if err != nil {
		log.Errorf("[initlink] parse ip63 err: %v", err)
		return ip.From, 0, err
	}
	err = h.Unmarshal(ip.Data)
	if err != nil {
		log.Errorf("[initlink] parse hello err: %v", err)
		return ip.From, 0, err
	}
	delay := h.Time - t
	if delay <= 0 {
		log.Errorf("[initlink] tr: %v, t: %v", h.Time, t)
		return ip.From, delay, err
	}
	saveMap(ip.From, conn, mt)
	router.AddItem(ip.From, ip.From, uint16(delay/1000000))
	log.Printf("[initlink] %s 链接测试成功，延时%vns，对方ip: %x", url, delay, ip.From)
	return ip.From, delay, nil
}

var upgrader = websocket.Upgrader{}

func InitEntry(ent string) {
	myhello = hello.Hello{
		Entry: ent,
	}
	http.HandleFunc("/"+ent, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err == nil {
			go listen(conn)
		}
	})
}
