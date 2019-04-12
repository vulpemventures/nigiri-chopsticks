package router

import (
	"testing"

	cfg "github.com/vulpemventures/nigiri-chopsticks/config"
)

func TestNewRouter(t *testing.T) {
	router := newTestingRouter(false)
	if router == nil {
		t.Fatal("Failed to create router")
	}
	if router.RPCClient == nil {
		t.Fatal("Failed to create router's RPC client")
	}
	if router.Faucet == nil {
		t.Fatal("Failed to create router's faucet")
	}
	if router.Config.Chain() != "bitcoin" {
		t.Fatal("Router is not configured for bitcoin chain")
	}
}

func TestNewLiquidRouter(t *testing.T) {
	router := newTestingRouter(true)
	if router == nil {
		t.Fatal("Failed to create router")
	}
	if router.RPCClient == nil {
		t.Fatal("Failed to create router's RPC client")
	}
	if router.Faucet == nil {
		t.Fatal("Failed to create router's faucet")
	}
	if router.Config == nil {
		t.Fatal("Failed to create router's configuration")
	}
	if router.Config.Chain() != "liquid" {
		t.Fatal("Router is not configured for liquid sidechain")
	}
}

func newTestingRouter(isLiquid bool) *Router {
	config := cfg.NewTestConfig()
	if isLiquid {
		config = cfg.NewTestLiquidConfig()
	}

	return NewRouter(config)
}
