package ip64

import (
	"time"

	"github.com/sirupsen/logrus"
)

func (ip *Ip64) Pack(from uint64, to uint64, data *[]byte, prototype uint32) {
	ip.Prototype = prototype
	ip.From = from
	ip.To = to
	ip.Data = *data
	ip.Ttl = 8
	ip.Time = time.Now().UnixNano()
	logrus.Infof("[ip64] pack data: %d bytes.", len(ip.Data))
}
