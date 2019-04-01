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
		"tls_enabled":    config.IsTLSEnabled(),
		"faucet_enabled": config.IsFaucetEnabled(),
		"mining_enabled": config.IsMiningEnabled(),
		"logger_enabled": config.IsLoggerEnabled(),
		"listen_url":     config.ListenURL(),
		"electrs_url":    config.ElectrsURL(),
		"rpc_server_url": config.RPCServerURL(),
	}).Info("Starting server with configuration:")

	r := router.NewRouter(config)

	if !config.IsTLSEnabled() {
		s := makeHTTPServer(r)
		s.Addr = fmt.Sprintf(config.ListenURL())
		if err = s.ListenAndServe(); err != nil {
			log.WithError(err).Fatal("HTTP server exited with error")
		}
	}

	if config.IsTLSEnabled() {
		dataDir := "."
		m := &autocert.Manager{
			Prompt: autocert.AcceptTOS,
			Cache:  autocert.DirCache(dataDir),
		}

		s := makeHTTPServer(r)
		s.Addr = fmt.Sprintf(config.ListenURL())
		s.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}

		if err = s.ListenAndServeTLS("", ""); err != nil {
			log.WithError(err).Fatal("HTTPS server exited with error")
		}
	}
}
