package link

import (
	"fmt"
	"net/http"
	"time"

	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/ip64"
	"github.com/IPoWS/node-core/router"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

var (
	npsurl  string
	Mywsip  uint64
	myhello hello.Hello
)

func initLink(conn *websocket.Conn, adviceip uint64) (uint64, int64, error) {
	t := time.Now().UnixNano()
	h := myhello
	if adviceip > 0 {
		h.Mask = 0xffff_ffff_0000_0000
	}
	sendHelloUnknown(conn, &h, adviceip)
	_, p, err := conn.ReadMessage()
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return adviceip, 0, err
	}
	var ip ip64.Ip64
	err = ip.Unmarshal(p)
	if err != nil {
		log.Errorf("[initlink] parse ip64 err: %v", err)
		return ip.From, 0, err
	}
	err = h.Unmarshal(ip.Data)
	if err != nil {
		log.Errorf("[initlink] parse hello err: %v", err)
		return ip.From, 0, err
	}
	delay := ip.Time - t
	if delay <= 0 {
		log.Errorf("[initlink] tr: %v, t: %v", ip.Time, t)
		return ip.From, delay, err
	}
	if adviceip > 0 && ip.From&h.Mask != adviceip&h.Mask {
		log.Infof("[initlink] peer %x reported a diff wsip than adv %x.", ip.From, adviceip)
		return ip.From, delay, fmt.Errorf("peer %x reported a diff wsip than adv %x.", ip.From, adviceip)
	}
	if ip.From > 0 {
		saveMap(ip.From, conn)
		router.AddItem(ip.From, ip.From, uint16(delay/100000))
		NodesList.AddNode(h.Host, h.Entry, ip.From, h.Name, uint64(delay))
		registerNode(ip.From)
		SetAlive(ip.From)
		log.Printf("[initlink] 链接测试成功，延时%vns，对方ip: %x", delay, ip.From)
		return ip.From, delay, nil
	} else {
		saveMap(ip.To, conn)
	}
	log.Printf("[initlink] 链接测试失败，延时%vns，对方ip: %x", delay, ip.From)
	return ip.From, delay, fmt.Errorf("peer wsip is 0.")
}

// InitLink 初始化连接 返回 wsip delay error, url 必须以 ws:// 开头, 以 ent 结尾, adviceip 可为 0
func InitLink(url string, adviceip uint64) (uint64, int64, error) {
	log.Printf("[initlink] connecting to %s", url)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Errorf("[initlink] %v", err)
		return 0, 0, err
	}
	return initLink(conn, adviceip)
}

var upgrader = websocket.Upgrader{}

// InitEntry 初始化 ws entry, nps 形如 ws://xxx/nps
func InitEntry(nps string, ent string, hostname string, mask uint64) {
	npsurl = nps
	myhello = hello.Hello{
		Entry: ent,
		Name:  hostname,
		Mask:  mask,
	}
	logrus.Infoln("[InitEntry] ", nps)
	initEntry(ent)
}

func initEntry(ent string) {
	http.HandleFunc("/"+ent, func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err == nil {
			go listen(conn)
		}
	})
}

// UpgradeLink 直接从 http 请求升级连接
func UpgradeLink(w http.ResponseWriter, r *http.Request, adviceip uint64) (uint64, int64, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err == nil {
		return initLink(conn, adviceip)
	}
	return 0, 0, err
}
