package router

import (
	"github.com/gorilla/mux"
	cfg "github.com/vulpemventures/nigiri-chopsticks/config"
	"github.com/vulpemventures/nigiri-chopsticks/faucet"
	"github.com/vulpemventures/nigiri-chopsticks/faucet/regtest"
	"github.com/vulpemventures/nigiri-chopsticks/router/middleware"
)

// Router extends gorilla Router
type Router struct {
	*mux.Router
	Config *cfg.Config
	Faucet faucet.Faucet
}

// NewRouter returns a new Router instance
func NewRouter(config *cfg.Config) *Router {
	router := mux.NewRouter().StrictSlash(true)

	r := &Router{router, config, nil}

	if config.Server.FaucetEnabled {
		url := r.Config.RPCServerURL()
		r.Faucet = regtestfaucet.NewFaucet(url)
		r.HandleFunc("/faucet", r.HandleFaucetRequest).Methods("POST")
	}

	r.Use(middleware.Logger)
	r.HandleFunc("/tx", r.ProxyBroadcast).Methods("POST")
	r.PathPrefix("/").HandlerFunc(r.ProxyElectrs)

	return r
}
