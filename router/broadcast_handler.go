package router

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// HandleBroadcastRequest forwards the request to the electrs HTTP server and mines a block if mining is enabled
func (r *Router) HandleBroadcastRequest(res http.ResponseWriter, req *http.Request) {
	r.HandleElectrsRequest(res, req)

	if r.Config.IsMiningEnabled() {
		blocks := 1
		if r.Config.Chain() == "liquid" {
			blocks = 10
		}
		status, blockHashes, err := r.Faucet.Mine(blocks)
		if err != nil {
			log.WithError(err).WithField("status", status).Warning("An unexpected error occured while mining blocks")
		} else {
			if r.Config.IsLoggerEnabled() {
				log.WithField("blocks mined", len(blockHashes)).Info("Transaction has been confirmed")
			}
		}
	}
}
