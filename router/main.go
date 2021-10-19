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

	// here we are forcing always the calls against the bitcoin/elements default wallet ""
	rpcClient, _ := helpers.NewRpcClient(config.RPCServerURL()+"/wallet/", false, 10)

	r := &Router{router, config, rpcClient, nil, nil}

	// Handle all preflight request
	r.Router.Methods(http.MethodOptions).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// fmt.Printf("OPTIONS")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.WriteHeader(http.StatusNoContent)
	})

	// From Bitcoin core 0.21 the default wallet "" is not created anymore.
	//So we check if none is already loaded and we create it
	err := helpers.CreateWalletIfNotExists(rpcClient)
	if err != nil {
		log.WithError(err).Fatalln("creating wallet")
	}
	log.Debug("empty wallet has been created")

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

		// From Elements core 0.21 if we use initialfreecoins we must rescan the chain
		err = helpers.RescanBlockchain(rpcClient)
		if err != nil {
			log.WithError(err).Fatalln("rescan blockchain")
		}
		log.Debug("rescan completed")
	}

	if config.IsLoggerEnabled() {
		r.Use(middleware.Logger)
	}

	r.HandleFunc("/getnewaddress", r.HandleAddressRequest).Methods(http.MethodGet, http.MethodOptions)
	r.HandleFunc("/tx", r.HandleBroadcastRequest).Methods(http.MethodPost, http.MethodOptions)
	r.PathPrefix("/").HandlerFunc(r.HandleElectrsRequest)

	return r
}
