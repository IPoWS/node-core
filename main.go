package main

import (
	"fmt"

	"github.com/IPoWS/node-core/link"
)

func main() {
	link.SetNPSUrl("http://127.0.0.1:8080")
	link.InitEntry("123456")
	fmt.Print(link.RegisterAndGetNodes("123456"))
	select {}
}
