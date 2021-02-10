package router

import (
	"net/http"
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
	Registry  *helpers.Registry
}

// NewRouter returns a new Router instance
func NewRouter(config cfg.Config) *Router {
	router := mux.NewRouter().StrictSlash(true)
	router.Use(mux.CORSMethodMiddleware(router))

	rpcClient, _ := helpers.NewRpcClient(config.RPCServerURL(), false, 10)

	r := &Router{router, config, rpcClient, nil, nil}

	if r.Config.IsFaucetEnabled() {
		faucet := faucet.NewFaucet(config.RPCServerURL(), rpcClient)
		r.Faucet = faucet
		r.HandleFunc("/faucet", r.HandleFaucetRequest).Methods(http.MethodPost, http.MethodOptions)
		if config.Chain() == "liquid" {
			registry, _ := helpers.NewRegistry(config.RegistryPath())
			r.Registry = registry
			r.HandleFunc("/mint", r.HandleMintRequest).Methods(http.MethodPost, http.MethodOptions)
			r.HandleFunc("/registry", r.HandleRegistryRequest).Methods(http.MethodPost, http.MethodOptions)
		}

		var numBlockToGenerate int = 1
		if config.Chain() == "bitcoin" {
			numBlockToGenerate = 101
		}
		status, blockHashes, err := r.Faucet.Fund(numBlockToGenerate)

		for err != nil && strings.Contains(err.Error(), "Loading") && status == 500 {
			time.Sleep(2 * time.Second)
			status, blockHashes, err = r.Faucet.Fund(numBlockToGenerate)
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

	r.HandleFunc("/address", r.HandleAddressRequest).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/tx", r.HandleBroadcastRequest).Methods(http.MethodPost, http.MethodOptions)
	r.PathPrefix("/").HandlerFunc(r.HandleElectrsRequest).Methods(http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodOptions)

	return r
}
