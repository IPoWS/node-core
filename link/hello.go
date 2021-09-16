package link

import (
	"github.com/IPoWS/node-core/data/hello"
	"github.com/IPoWS/node-core/ip64"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
)

// SendHello 从自身发送 hello 给对方，链接要靠 sendhello 保持交互，同时注册下一个 listen 函数
func SendHello(to uint64, conn *websocket.Conn) ([]byte, error) {
	h := myhello
	return sendHello(to, &h, listen, conn)
}

// sendHello 发送 hello 给对方
func sendHello(wsip uint64, h *hello.Hello, handler func(*websocket.Conn), conn *websocket.Conn) ([]byte, error) {
	data, err := h.Marshal()
	if err == nil {
		var ip ip64.Ip64
		ip.Pack(Mywsip, wsip, &data, ip64.HelloType, 0, 0)
		logrus.Infof("[sendHello] send hello to %x.", wsip)
		data, err = ip.Send(conn, websocket.BinaryMessage, handler)
	}
	if err != nil {
		logrus.Errorf("[sendHello] %v", err)
	}
	return data, err
}

// sendHelloUnknown 发送 hello 给未知 ip 方
func sendHelloUnknown(conn *websocket.Conn, h *hello.Hello, adviceip uint64) ([]byte, error) {
	data, err := h.Marshal()
	if err == nil {
		var ip ip64.Ip64
		ip.Pack(Mywsip, adviceip, &data, ip64.HelloType, 0, 0)
		logrus.Info("[sendHello] send hello.")
		data, err = ip.Send(conn, websocket.BinaryMessage, nil)
	}
	if err != nil {
		logrus.Errorf("[sendHello] %v", err)
	}
	return data, err
}
