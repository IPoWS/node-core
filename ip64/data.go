package ip64

import (
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func (ip *Ip64) Send(conn *websocket.Conn, mt int) error {
	ip.Ttl -= 1
	if ip.Ttl <= 0 {
		return fmt.Errorf("[ip64] send to %x failed: ttl = %d.", ip.To, ip.Ttl)
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
