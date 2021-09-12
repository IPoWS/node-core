package main

import (
	"github.com/IPoWS/node-core/link"
)

func main() {
	link.SetNPSUrl("127.0.0.1:8080")
	link.InitEntry("123456")
	select {}
}
