package router

import (
	"testing"

	"github.com/vulpemventures/nigiri-chopsticks/config"
)

func NewTestRouter(liquid bool) *Router {
	if liquid {
		return NewRouter(config.NewLiquidTestConfig())
	}

	return NewRouter(config.NewTestConfig())
}

func TestBitcoinChopsticks(t *testing.T) {
	r := NewTestRouter(false)
	if r == nil {
		t.Fatal("Expected *Router, got <nil>")
	}
}

func TestLiquidChopsticks(t *testing.T) {
	r := NewTestRouter(true)
	if r == nil {
		t.Fatal("Expected *Router, got <nil>")
	}
}
