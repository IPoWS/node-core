package link

import (
	"fmt"

	"github.com/IPoWS/node-core/ip64"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

func Send(to uint64, data *[]byte, prototype uint32, port uint16) error {
	sendmu.RLock()
	wsn, ok := sendmap[to]
	sendmu.RUnlock()
	if ok {
		var ip ip64.Ip64
		ip.Pack(Mywsip, to, data, prototype|(uint32(port)<<16))
		logrus.Infof("[Send] link send %d bytes to %x.", len(*data), to)
		return ip.Send(wsn, websocket.BinaryMessage)
	}
	return fmt.Errorf("dest %x unreachable.", to)
}

func Forward(to uint64, ip *ip64.Ip64) error {
	sendmu.RLock()
	wsn, ok := sendmap[to]
	sendmu.RUnlock()
	if ok {
		return ip.Send(wsn, websocket.BinaryMessage)
	}
	return fmt.Errorf("dest %x unreachable.", to)
}
