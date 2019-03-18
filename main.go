package main

import (
	"fmt"
	"net/http"

	cfg "github.com/vulpemventures/nigiri-chopsticks/config"
	"github.com/vulpemventures/nigiri-chopsticks/router"

	log "github.com/sirupsen/logrus"
)

func main() {
	config, err := cfg.NewConfigFromFlags()
	if err != nil {
		log.WithError(err).Fatal("Failed to parse flags")
	}

	log.WithFields(log.Fields{
		"address":         fmt.Sprintf("%s://%s:%s", config.Server.Proto, config.Server.Host, config.Server.Port),
		"electrs_address": fmt.Sprintf("%s:%s", config.Electrs.Host, config.Electrs.Port),
		"bitcoin_address": fmt.Sprintf("%s:%s", config.Bitcoin.Host, config.Bitcoin.Port),
		"bitcoin_cookie":  fmt.Sprintf("%s:%s", config.Bitcoin.RPCUser, config.Bitcoin.RPCPassword),
		"liquid_address":  fmt.Sprintf("%s:%s", config.Liquid.Host, config.Liquid.Port),
	}).Info("Starting server with configuration details:")

	r := router.NewRouter(config)

	if config.Server.Proto == "http" {
		if err = http.ListenAndServe(fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port), r); err != nil {
			log.Fatal(err)
		}
	}
}
