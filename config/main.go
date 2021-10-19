package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultTLSEnabled    = false
	defaultLoggerEnabled = false
	defaultFaucetEnabled = false
	defaultMiningEnabled = false

	defaultAddr        = "localhost:3000"
	defaultElectrsAddr = "localhost:3002"
	defaultRPCAddr     = "localhost:19001"
	defaultRPCCookie   = "admin1:123"
	defaultChain       = "bitcoin"
)

var defaultRegistryPath, _ = os.Getwd()

// Config type is used to parse flag options
type Config interface {
	IsTLSEnabled() bool
	IsFaucetEnabled() bool
	IsLoggerEnabled() bool
	IsMiningEnabled() bool
	ListenURL() string
	RPCServerURL() string
	ElectrsURL() string
	Chain() string
	RegistryPath() string
}

type config struct {
	server struct {
		tlsEnabled    bool
		faucetEnabled bool
		miningEnabled bool
		loggerEnabled bool
		host          string
		port          string
		chain         string
		registryPath  string
	}
	electrs struct {
		host string
		port string
	}
	rpcServer struct {
		user     string
		password string
		host     string
		port     string
	}
}

// NewConfigFromFlags parses flags and returns a Config
func NewConfigFromFlags() (Config, error) {
	tlsEnabled := flag.Bool("use-tls", defaultTLSEnabled, "Set true to use https")
	faucetEnabled := flag.Bool("use-faucet", defaultFaucetEnabled, "Set to use faucet")
	miningEnabled := flag.Bool("use-mining", defaultMiningEnabled, "Set to false to disable block mining right after broadcasting requests")
	loggerEnabled := flag.Bool("use-logger", defaultLoggerEnabled, "Set true to log every request/response")

	addr := flag.String("addr", defaultAddr, "Chopsticks listen address")
	electrsAddr := flag.String("electrs-addr", defaultElectrsAddr, "Elctrs HTTP server address")
	rpcAddr := flag.String("rpc-addr", defaultRPCAddr, "RPC server address")
	rpcCookie := flag.String("rpc-cookie", defaultRPCCookie, "RPC server user and password")
	chain := flag.String("chain", defaultChain, "Set default chain. Eihter 'bitcoin' or 'liquid'")
	registryPath := flag.String("registry-path", defaultRegistryPath, "(Liquid only) Set path for asset registry JSON file")
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

	if !filepath.IsAbs(*registryPath) {
		flag.PrintDefaults()
		return nil, fmt.Errorf("Invalid registry path")
	}

	c := &config{}
	c.server.loggerEnabled = *loggerEnabled
	c.server.tlsEnabled = *tlsEnabled
	c.server.faucetEnabled = *faucetEnabled
	c.server.miningEnabled = *miningEnabled
	c.server.host = host
	c.server.port = port
	c.server.chain = *chain
	c.server.registryPath = *registryPath

	c.electrs.host = electrsHost
	c.electrs.port = electrsPort

	c.rpcServer.host = rpcHost
	c.rpcServer.port = rpcPort
	c.rpcServer.user = rpcUser
	c.rpcServer.password = rpcPassword

	return c, nil
}

func (c *config) IsTLSEnabled() bool {
	return c.server.tlsEnabled
}

func (c *config) IsFaucetEnabled() bool {
	return c.server.faucetEnabled
}

func (c *config) IsLoggerEnabled() bool {
	return c.server.loggerEnabled
}

func (c *config) IsMiningEnabled() bool {
	return c.server.miningEnabled
}

func (c *config) ListenURL() string {
	return fmt.Sprintf("%s:%s", c.server.host, c.server.port)
}

func (c *config) RPCServerURL() string {
	return fmt.Sprintf("http://%s:%s@%s:%s", c.rpcServer.user, c.rpcServer.password, c.rpcServer.host, c.rpcServer.port)
}

func (c *config) ElectrsURL() string {
	return fmt.Sprintf("http://%s:%s", c.electrs.host, c.electrs.port)
}

func (c *config) Chain() string {
	return c.server.chain
}

func (c *config) RegistryPath() string {
	return c.server.registryPath
}

func splitString(addr string) (string, string, bool) {
	if splitAddr := strings.Split(addr, ":"); len(splitAddr) == 2 {
		return splitAddr[0], splitAddr[1], true
	}

	return "", "", false
}

func NewTestConfig() Config {
	c := &config{}
	rpcPort := "18433"
	if os.Getenv("CI") == "true" {
		rpcPort = "18443"
	}

	c.server.tlsEnabled = false
	c.server.loggerEnabled = false
	c.server.faucetEnabled = true
	c.server.miningEnabled = true
	c.server.host = "localhost"
	c.server.port = "7000"
	c.server.chain = "bitcoin"

	c.electrs.host = "localhost"
	c.electrs.port = "30000"

	c.rpcServer.host = "localhost"
	c.rpcServer.port = rpcPort
	c.rpcServer.user = "admin1"
	c.rpcServer.password = "123"

	return c
}

func NewLiquidTestConfig() Config {
	c := &config{}
	rpcPort := "18884"
	if os.Getenv("CI") == "true" {
		rpcPort = "18884"
	}

	c.server.tlsEnabled = false
	c.server.loggerEnabled = false
	c.server.faucetEnabled = true
	c.server.miningEnabled = true
	c.server.host = "localhost"
	c.server.port = "7001"
	c.server.chain = "liquid"
	c.server.registryPath = defaultRegistryPath

	c.electrs.host = "localhost"
	c.electrs.port = "30001"

	c.rpcServer.host = "localhost"
	c.rpcServer.port = rpcPort
	c.rpcServer.user = "admin1"
	c.rpcServer.password = "123"

	return c
}
