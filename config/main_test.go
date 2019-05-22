package config

import "testing"

func TestNewConfigs(t *testing.T) {
	conf := NewTestConfig()
	if conf == nil {
		t.Fatal("Expected Config interface, got <nil>")
	}

	conf = NewLiquidTestConfig()
	if conf == nil {
		t.Fatal("Expected Config interface, got <nil>")
	}
}
