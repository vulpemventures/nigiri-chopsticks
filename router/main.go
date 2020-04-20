package router

import (
	"strings"
	"time"

	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	cfg "github.com/vulpemventures/nigiri-chopsticks/config"
	"github.com/vulpemventures/nigiri-chopsticks/faucet"
	"github.com/vulpemventures/nigiri-chopsticks/helpers"
	"github.com/vulpemventures/nigiri-chopsticks/router/middleware"
)

// Router extends gorilla Router
type Router struct {
	*mux.Router
	Config    cfg.Config
	RPCClient *helpers.RpcClient
	Faucet    *faucet.Faucet
}

// NewRouter returns a new Router instance
func NewRouter(config cfg.Config) *Router {
	router := mux.NewRouter().StrictSlash(true)
	rpcClient, _ := helpers.NewRpcClient(config.RPCServerURL(), false, 10)

	r := &Router{router, config, rpcClient, nil}

	if r.Config.IsFaucetEnabled() {
		faucet := faucet.NewFaucet(config.RPCServerURL(), rpcClient)
		r.Faucet = faucet
		r.HandleFunc("/faucet", r.HandleFaucetRequest).Methods("POST")
		if config.Chain() == "liquid" {
			r.HandleFunc("/mint", r.HandleMintRequest).Methods("POST")
		}

		status, blockHashes, err := r.Faucet.Fund()
		for err != nil && strings.Contains(err.Error(), "Loading") && status == 500 {
			time.Sleep(2 * time.Second)
			status, blockHashes, err = r.Faucet.Fund()
		}
		if err != nil {
			log.WithField("status", status).WithError(err).Warning("Faucet not funded, check the error")
		}
		if len(blockHashes) > 0 {
			log.WithField("num_blocks", len(blockHashes)).Info("Faucet has been funded mining some blocks")
		}
	}

	if config.IsLoggerEnabled() {
		r.Use(middleware.Logger)
	}
	r.HandleFunc("/tx", r.HandleBroadcastRequest).Methods("POST")
	r.PathPrefix("/").HandlerFunc(r.HandleElectrsRequest)

	return r
}
