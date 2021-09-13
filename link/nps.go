package link

import (
	"io"
	"net/http"

	"github.com/IPoWS/node-core/router"
	"github.com/sirupsen/logrus"
)

func Register(ent string) {
	resp, err := http.Get(npsurl + "?ent=" + ent)
	if err == nil {
		data, err := io.ReadAll(resp.Body)
		if err == nil {
			router.ParseRawNodes(data)
		}
	}
	if err != nil {
		logrus.Errorf("[RegisterAndGetNodes] %v", err)
	}
}
