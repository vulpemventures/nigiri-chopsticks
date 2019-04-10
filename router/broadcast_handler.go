package router

import (
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleBroadcastRequest forwards the request to the electrs HTTP server and mines a block if mining is enabled
func (r *Router) HandleBroadcastRequest(res http.ResponseWriter, req *http.Request) {
	/*
		since our electrs version does not expose a broadcast endpoint for liquid, we need
		to call the sendrawtransaction method of the underlying RPC daaemon instead
	*/
	if r.Config.Chain() == "liquid" {
		tx := req.URL.Query().Get("tx")

		status, resp, err := r.RPCClient.Call("sendrawtransaction", []interface{}{tx})
		if err != nil {
			http.Error(res, err.Error(), status)
			return
		}

		var txId string
		if err = json.Unmarshal(resp.Result, &txId); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		json.NewEncoder(res).Encode(map[string]string{"txId": txId})
	} else {
		r.HandleElectrsRequest(res, req)
	}

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
