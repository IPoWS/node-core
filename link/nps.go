package link

import (
	"io"
	"net/http"

	"github.com/IPoWS/node-core/data/nodes"
)

func GetNodes(ent string) *nodes.Nodes {
	var nodesList nodes.Nodes
	resp, err := http.Get(npsurl + "?ent=" + ent)
	if err == nil {
		data, err := io.ReadAll(resp.Body)
		if err == nil {
			err = nodesList.Unmarshal(data)
			if err == nil {
				return &nodesList
			}
		}
	}
	return nil
}
