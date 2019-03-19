package config

import (
	"flag"
	"fmt"
	"strings"
)

const (
	defaultTLSEnabled    = false
	defaultFaucetEnabled = true

	defaultAddr        = "localhost:3000"
	defaultElectrsAddr = "localhost:3002"
	defaultFaucetAddr  = "localhost:3001"

	defaultBtcCookie = "admin1:123"
)

// Config type is used to parse flag options
type Config struct {
	Server struct {
		TLSEnabled    bool
		FaucetEnabled bool
		Host          string
		Port          string
	}
	Electrs struct {
		Host string
		Port string
	}
	Faucet struct {
		Host string
		Port string
	}
}

// NewConfigFromFlags parses flags and returns a Config
func NewConfigFromFlags() (Config, error) {
	tlsEnabled := flag.Bool("use-tls", defaultTLSEnabled, "Set true to use https}")
	faucetEnabled := flag.Bool("use-faucet", defaultFaucetEnabled, "Set to true to use faucet")

	addr := flag.String("addr", defaultAddr, "Listen address")
	electrsAddr := flag.String("electrs-addr", defaultElectrsAddr, "Elctrs HTTP server address")
	faucetAddr := flag.String("faucet-addr", defaultFaucetAddr, "Faucet server address")
	flag.Parse()

	config := Config{}

	host, port, ok := splitString(*addr)
	if !ok {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid server address")
	}

	electrsHost, electrsPort, ok := splitString(*electrsAddr)
	if !ok {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid electrs HTTP server address")
	}

	faucetHost, faucetPort, ok := splitString(*faucetAddr)
	if !ok {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid faucet HTTP server address")
	}

	c := Config{}
	c.Server.TLSEnabled = *tlsEnabled
	c.Server.FaucetEnabled = *faucetEnabled
	c.Server.Host = host
	c.Server.Port = port

	c.Electrs.Host = electrsHost
	c.Electrs.Port = electrsPort

	c.Faucet.Host = faucetHost
	c.Faucet.Port = faucetPort

	return c, nil
}

func splitString(addr string) (string, string, bool) {
	if splitAddr := strings.Split(addr, ":"); len(splitAddr) == 2 {
		return splitAddr[0], splitAddr[1], true
	}

	return "", "", false
}
