package router

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleBroadcastRequest forwards the request to the electrs HTTP server and mines a block if mining is enabled
func (r *Router) HandleBroadcastRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")
	r.HandleElectrsRequest(res, req)

	if r.Config.IsMiningEnabled() {
		status, blockHashes, err := r.Faucet.Mine(1)
		if err != nil {
			log.WithError(err).WithField("status", status).Warning("An unexpected error occured while mining blocks")
		} else {
			if r.Config.IsLoggerEnabled() {
				log.WithField("blocks mined", len(blockHashes)).Info("Transaction has been confirmed")
			}
		}
	}
}
