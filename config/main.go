package config

import (
	"flag"
	"fmt"
	"strings"
)

const (
	defaultProto = "http"

	defaultAddr        = "localhost:3000"
	defaultElectrsAddr = "localhost:3002"
	defaultBtcAddr     = "localhost:19001"
	defaultLiquidAddr  = "localhost:18884"

	defaultBtcCookie = "admin1:123"
)

// Config type is used to parse flag options
type Config struct {
	Server struct {
		Proto string
		Host  string
		Port  string
	}
	Electrs struct {
		Host string
		Port string
	}
	Bitcoin struct {
		Host        string
		Port        string
		RPCUser     string
		RPCPassword string
	}
	Liquid struct {
		Host string
		Port string
	}
}

// NewConfigFromFlags parses flags and returns a Config
func NewConfigFromFlags() (Config, error) {
	proto := flag.String("proto", defaultProto, "Proto {http|https}")
	addr := flag.String("addr", defaultAddr, "Listen address")
	electrsAddr := flag.String("electrs-addr", defaultElectrsAddr, "Elctrs HTTP server address")
	btcAddr := flag.String("btc-addr", defaultBtcAddr, "Bitcoin RPC server address")
	btcCookie := flag.String("btc-cookie", defaultBtcCookie, "Colon separated (:) bitcoin RPC user and password")
	liquidAddr := flag.String("liquid-addr", defaultLiquidAddr, "Liquid RPC server address")
	flag.Parse()

	config := Config{}

	if *proto != "http" && *proto != "https" {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid proto")
	}

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

	btcHost, btcPort, ok := splitString(*btcAddr)
	if !ok {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid bitcoin RPC server address")
	}

	liquidHost, liquidPort, ok := splitString(*liquidAddr)
	if !ok {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid liquid RPC server address")
	}

	btcUser, btcPassword, ok := splitString(*btcCookie)
	if !ok {
		flag.PrintDefaults()
		return config, fmt.Errorf("Invalid bitcoin RPC cookie")
	}

	c := Config{}
	c.Server.Proto = *proto
	c.Server.Host = host
	c.Server.Port = port

	c.Electrs.Host = electrsHost
	c.Electrs.Port = electrsPort

	c.Bitcoin.Host = btcHost
	c.Bitcoin.Port = btcPort
	c.Bitcoin.RPCUser = btcUser
	c.Bitcoin.RPCPassword = btcPassword

	c.Liquid.Host = liquidHost
	c.Liquid.Port = liquidPort

	return c, nil
}

func splitString(addr string) (string, string, bool) {
	if splitAddr := strings.Split(addr, ":"); len(splitAddr) == 2 {
		return splitAddr[0], splitAddr[1], true
	}

	return "", "", false
}
