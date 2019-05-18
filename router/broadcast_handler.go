package router

import (
	"net/http"
)

// HandleBroadcastRequest forwards the request to the electrs HTTP server and mines a block if mining is enabled
func (r *Router) HandleBroadcastRequest(res http.ResponseWriter, req *http.Request) {
	r.HandleElectrsRequest(res, req)

	if r.Config.IsMiningEnabled() {
		r.Faucet.Mine(1)
	}
}
