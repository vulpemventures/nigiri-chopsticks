package router

import (
	"github.com/gorilla/mux"
	cfg "github.com/vulpemventures/nigiri-chopsticks/config"
)

// Router extends gorilla Router
type Router struct {
	*mux.Router
	Config cfg.Config
}

// NewRouter returns a new Router instance
func NewRouter(config cfg.Config) *Router {
	router := mux.NewRouter().StrictSlash(true)

	r := &Router{router, config}

	if config.Server.FaucetEnabled {
		r.HandleFunc("/faucet", r.ProxyFaucet).Methods("POST")
	}
	r.HandleFunc("/tx", r.ProxyBroadcast).Methods("POST")
	r.PathPrefix("/").HandlerFunc(r.ProxyElectrs)

	return r
}
