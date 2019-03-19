package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
	cfg "github.com/vulpemventures/nigiri-chopsticks/config"
	"github.com/vulpemventures/nigiri-chopsticks/router"
	"golang.org/x/crypto/acme/autocert"
)

func makeHTTPServer(handler http.Handler) *http.Server {
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      handler,
	}
}

func main() {
	config, err := cfg.NewConfigFromFlags()
	if err != nil {
		log.WithError(err).Fatal("Failed to parse flags")
	}

	log.WithFields(log.Fields{
		"tls-enabled":     config.Server.TLSEnabled,
		"faucet-enabled":  config.Server.FaucetEnabled,
		"mining-enabled":  config.Server.MiningEnabled,
		"address":         fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port),
		"electrs_address": fmt.Sprintf("%s:%s", config.Electrs.Host, config.Electrs.Port),
		"faucet_address":  fmt.Sprintf("%s:%s", config.Faucet.Host, config.Faucet.Port),
	}).Info("Starting server with configuration:")

	r := router.NewRouter(config)

	if !config.Server.TLSEnabled {
		s := makeHTTPServer(r)
		s.Addr = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
		if err = s.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}

	if config.Server.TLSEnabled {
		dataDir := "."
		m := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(dataDir),
		}

		s := makeHTTPServer(r)
		s.Addr = fmt.Sprintf("%s:%s", config.Server.Host, config.Server.Port)
		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		if err = s.ListenAndServeTLS("", ""); err != nil {
			log.WithError(err).Fatal("HTTPS server exited with error")
		}
	}
}
