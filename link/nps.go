package link

import "net/http"

func GetNodes(ent string) {
	http.Get(npsurl + "?ent=" + ent)
}
