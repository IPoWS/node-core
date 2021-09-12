package link

import (
	"io"
	"net/http"

	"github.com/IPoWS/node-core/data/nodes"
	"github.com/sirupsen/logrus"
)

func RegisterAndGetNodes(ent string) *nodes.Nodes {
	var nodesList = new(nodes.Nodes)
	resp, err := http.Get(npsurl + "?ent=" + ent)
	if err == nil {
		data, err := io.ReadAll(resp.Body)
		if err == nil {
			err = nodesList.Unmarshal(data)
			if err == nil {
				return nodesList
			}
		}
	}
	if err != nil {
		logrus.Errorf("[RegisterAndGetNodes] %v", err)
	}
	return nil
}
