package main

import (
	"github.com/IPoWS/node-core/link"
)

func main() {
	link.SetNPSUrl("http://127.0.0.1:8080/nps")
	link.InitEntry("123456")
	link.Register("123456")
	select {}
}
