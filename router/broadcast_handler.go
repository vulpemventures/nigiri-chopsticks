package router

import (
	"encoding/json"
	"net/http"

	"github.com/vulpemventures/nigiri-chopsticks/helpers/liquidrpc"

	log "github.com/sirupsen/logrus"
	rpchelper "github.com/vulpemventures/nigiri-chopsticks/helpers/rpc"
)

// HandleBroadcastRequest forwards the request to the electrs HTTP server and mines a block if mining is enabled
func (r *Router) HandleBroadcastRequest(res http.ResponseWriter, req *http.Request) {
	/*
		since our electrs version does not expose a broadcast endpoint for liquid, we need
		to call the sendrawtransaction method of the underlying RPC daaemon instead
	*/
	if r.Config.Chain() == "liquid" {
		url := r.Config.RPCServerURL()
		tx := req.URL.Query().Get("tx")

		status, txID, err := liquidrpc.Sendrawtransaction(url, tx)
		if err != nil {
			http.Error(res, err.Error(), status)
			return
		}

		json.NewEncoder(res).Encode(map[string]string{"txId": txID})
	} else {
		r.HandleElectrsRequest(res, req)
	}

	if r.Config.IsMiningEnabled() {
		mineOneBlock(r)
	}
}

func mineOneBlock(r *Router) {
	url := r.Config.RPCServerURL()
	_, blockHash, err := rpchelper.Mine(url, 1)
	if r.Config.IsLoggerEnabled() {
		if err != nil {
			log.WithError(err).Warning("Error while mining a block")
		} else {
			log.WithField("block_hash", blockHash[0]).Info("Block mined")
		}
	}
}
