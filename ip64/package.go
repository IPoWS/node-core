package ip64

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (ip *Ip64) Pack(from uint64, to uint64, data *[]byte, proto uint16, srcport uint16, destport uint16) {
	ip.Destproto = uint32(proto) | (uint32(destport) << 16)
	ip.From = from
	ip.To = to
	ip.Data = *data
	ip.Srcttl = (int32(srcport) << 16) | 8
	ip.Time = time.Now().UnixNano()
	logrus.Debugf("[ip64] pack data: %d bytes.", len(ip.Data))
}
