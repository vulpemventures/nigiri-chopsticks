package router

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// ProxyBroadcast forwards tx broadcast request to electrs and mines a block if using chopsticks
func (r *Router) ProxyBroadcast(res http.ResponseWriter, req *http.Request) {
	r.ProxyElectrs(res, req)

	if r.Config.Server.MiningEnabled {
		url := fmt.Sprintf("http://%s:%s", r.Config.Faucet.Host, r.Config.Faucet.Port)
		status, resp, err := post(url, "", nil)
		if err != nil {
			log.WithError(err).Warning("Error while mining a block")
		}
		if status != http.StatusOK {
			log.WithFields(log.Fields{
				"response": resp,
				"status":   status,
			}).Warning("Error while mining a block")
		}
	}
}
