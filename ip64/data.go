package ip64

import "github.com/gorilla/websocket"

func (ip *Ip64) Send(conn *websocket.Conn, mt int) error {
	d, err := ip.Marshal()
	if err == nil {
		err = conn.WriteMessage(mt, d)
	}
	return err
}
