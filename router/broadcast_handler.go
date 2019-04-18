package router

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleBroadcastRequest forwards the request to the electrs HTTP server and mines a block if mining is enabled
func (r *Router) HandleBroadcastRequest(res http.ResponseWriter, req *http.Request) {
	r.HandleElectrsRequest(res, req)

	if r.Config.IsMiningEnabled() {
		mineOneBlock(r)
	}
}

func mineOneBlock(r *Router) {
	status, resp, err := r.RPCClient.Call("generate", []interface{}{1})
	if r.Config.IsLoggerEnabled() {
		if err != nil {
			log.WithError(err).WithField("status", status).Warning("Error while mining a block")
		} else {
			blockHashes := []string{}
			json.Unmarshal(resp.Result, &blockHashes)
			log.WithField("block_hash", blockHashes[0]).Info("Block mined")
		}
	}
}
