package ip64

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (ip *Ip64) Send(conn *websocket.Conn, mt int) error {
	ttl := ip.Srcttl & 0x0000_ffff
	ttl -= 1
	if ttl <= 0 {
		return fmt.Errorf("[ip64] send to %x failed: ttl = %d.", ip.To, ttl)
	}
	d, err := ip.Marshal()
	if err == nil {
		err = conn.WriteMessage(mt, d)
	}
	if err != nil {
		logrus.Errorln("[ip64] send err: ", err)
	}
	return err
}
