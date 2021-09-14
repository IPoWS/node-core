package ip64

import (
	"fmt"

	"github.com/gorilla/websocket"
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
	return err
}
