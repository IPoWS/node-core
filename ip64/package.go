package ip64

import (
	. "github.com/IPoWS/node-core/data"
)

func (ip *Ip64) Pack(from uint64, to uint64, fromport uint16, toport uint16, data interface{}, datalen uintptr, prototype uint32) {
	ip.Prototype = prototype
	ip.From = from
	ip.To = to
	ip.Ports = uint32(from)<<16 | uint32(to)
	ip.Datlen = uint32(datalen)
	ip.Data = Interface2Bytes(data, datalen)
}
