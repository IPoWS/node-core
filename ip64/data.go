package ip64

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// Send 发送 ip 包，如处理函数为空则直接返回数据
func (ip *Ip64) Send(conn *websocket.Conn, mt int, handleretdat func(data []byte, conn *websocket.Conn)) ([]byte, error) {
	ttl := ip.Srcttl & 0x0000_ffff
	ttl -= 1
	if ttl <= 0 {
		return nil, fmt.Errorf("[ip64] send to %x failed: ttl = %d.", ip.To, ttl)
	}
	d, err := ip.Marshal()
	if err == nil {
		err = conn.WriteMessage(mt, d)
	}
	if err == nil {
		_, d, err = conn.ReadMessage()
	}
	if err != nil {
		logrus.Errorln("[ip64] send err: ", err)
	} else if handleretdat != nil {
		handleretdat(d, conn)
	}
	return d, err
}
