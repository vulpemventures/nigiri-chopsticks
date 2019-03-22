package config

import (
	"flag"
	"fmt"
	"strings"
)

const (
	defaultTLSEnabled    = false
	defaultFaucetEnabled = true
	defaultMiningEnabled = true

	defaultAddr        = "localhost:3000"
	defaultElectrsAddr = "localhost:3002"
	defaultRPCAddr     = "localhost:19001"
	defaultRPCCookie   = "admin1:123"
)

// Config type is used to parse flag options
type Config struct {
	Server struct {
		TLSEnabled    bool
		FaucetEnabled bool
		MiningEnabled bool
		Host          string
		Port          string
	}
	Electrs struct {
		Host string
		Port string
	}
	RPCServer struct {
		User     string
		Password string
		Host     string
		Port     string
	}
}

// NewConfigFromFlags parses flags and returns a Config
func NewConfigFromFlags() (*Config, error) {
	tlsEnabled := flag.Bool("use-tls", defaultTLSEnabled, "Set true to use https}")
	faucetEnabled := flag.Bool("use-faucet", defaultFaucetEnabled, "Set to true to use faucet")
	miningEnabled := flag.Bool("use-mining", defaultMiningEnabled, "set to false to disable block mining right after broadcasting requests")

	addr := flag.String("addr", defaultAddr, "Listen address")
	electrsAddr := flag.String("electrs-addr", defaultElectrsAddr, "Elctrs HTTP server address")
	rpcAddr := flag.String("rpc-addr", defaultRPCAddr, "RPC server address")
	rpcCookie := flag.String("rpc-cookie", defaultRPCCookie, "RPC server user and password")
	flag.Parse()

	host, port, ok := splitString(*addr)
	if !ok {
		flag.PrintDefaults()
		return nil, fmt.Errorf("Invalid server address")
	}

	electrsHost, electrsPort, ok := splitString(*electrsAddr)
	if !ok {
		flag.PrintDefaults()
		return nil, fmt.Errorf("Invalid electrs HTTP server address")
	}

	rpcHost, rpcPort, ok := splitString(*rpcAddr)
	if !ok {
		flag.PrintDefaults()
		return nil, fmt.Errorf("Invalid RPC server address")
	}

	rpcUser, rpcPassword, ok := splitString(*rpcCookie)
	if !ok {
		flag.PrintDefaults()
		return nil, fmt.Errorf("Invalid RPC server cookie")
	}

	c := &Config{}
	c.Server.TLSEnabled = *tlsEnabled
	c.Server.FaucetEnabled = *faucetEnabled
	c.Server.MiningEnabled = *miningEnabled
	c.Server.Host = host
	c.Server.Port = port

	c.Electrs.Host = electrsHost
	c.Electrs.Port = electrsPort

	c.RPCServer.Host = rpcHost
	c.RPCServer.Port = rpcPort
	c.RPCServer.User = rpcUser
	c.RPCServer.Password = rpcPassword

	return c, nil
}

func (c *Config) RPCServerURL() string {
	return fmt.Sprintf("http://%s:%s@%s:%s", c.RPCServer.User, c.RPCServer.Password, c.RPCServer.Host, c.RPCServer.Port)
}

func (c *Config) ElectrsURL() string {
	return fmt.Sprintf("http://%s:%s", c.Electrs.Host, c.Electrs.Port)
}

func splitString(addr string) (string, string, bool) {
	if splitAddr := strings.Split(addr, ":"); len(splitAddr) == 2 {
		return splitAddr[0], splitAddr[1], true
	}

	return "", "", false
}
