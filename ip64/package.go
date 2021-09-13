package ip64

import (
	"github.com/sirupsen/logrus"
)

func (ip *Ip64) Pack(from uint64, to uint64, fromport uint16, toport uint16, data *[]byte, prototype uint32) {
	ip.Prototype = prototype
	ip.From = from
	ip.To = to
	ip.Ports = uint32(from)<<16 | uint32(to)
	ip.Data = *data
	logrus.Infof("[ip64] pack data: %d bytes.", len(ip.Data))
}
