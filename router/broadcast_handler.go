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
		url := r.Config.RPCServerURL()
		body := `{"jsonrpc":"1.0", "id": "2", "method":"generate", "params":[1]}`
		status, resp, err := post(url, body, nil)
		if err != nil {
			log.WithError(err).Warning("Error while mining a block")
		}
		if status != http.StatusOK {
			log.WithFields(log.Fields{
				"response": resp,
				"status":   status,
			}).Warning("Error while mining a block")
		} else {
			out := map[string]string{}
			json.Unmarshal([]byte(resp), &out)
			log.WithField("block_hash", out["result"]).Info("Block mined")
		}
	}
}
