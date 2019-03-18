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

	r.HandleFunc("/faucet/send", r.Send).Methods("POST")
	r.HandleFunc("/faucet/broadcast", r.Broadcast).Methods("POST")
	r.PathPrefix("/esplora/").HandlerFunc(r.ProxyElectrs)
	// r.PathPrefix("/regtest/").HandlerFunc(r.ProxyBitcoin)
	// r.PathPrefix("/liquid/").HandlerFunc(r.ProxyLiquid)

	return r
}
