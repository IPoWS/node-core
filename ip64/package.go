package ip64

import (
	. "github.com/IPoWS/node-core/data"
	"github.com/sirupsen/logrus"
)

func (ip *Ip64) Pack(from uint64, to uint64, fromport uint16, toport uint16, data interface{}, datalen uintptr, prototype uint32) {
	ip.Prototype = prototype
	ip.From = from
	ip.To = to
	ip.Ports = uint32(from)<<16 | uint32(to)
	ip.Data = Interface2Bytes(data, datalen)
	logrus.Infof("[ip64] pack data: %d bytes from original %d bytes.", len(ip.Data), datalen)
}
