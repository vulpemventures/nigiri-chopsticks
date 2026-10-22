package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/vulpemventures/nigiri-chopsticks/config"
)

// fakeNode is an in-memory JSON-RPC node that records every method it receives.
type fakeNode struct {
	*httptest.Server
	blockCount float64

	mu    sync.Mutex
	calls []rpcCall
}

type rpcCall struct {
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

func newFakeNode(t *testing.T, blockCount float64) *fakeNode {
	n := &fakeNode{blockCount: blockCount}
	n.Server = httptest.NewServer(http.HandlerFunc(n.handle))
	t.Cleanup(n.Close)
	return n
}

func (n *fakeNode) handle(w http.ResponseWriter, req *http.Request) {
	var call rpcCall
	if err := json.NewDecoder(req.Body).Decode(&call); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	n.mu.Lock()
	n.calls = append(n.calls, call)
	n.mu.Unlock()

	var result interface{}
	switch call.Method {
	case "listwallets":
		result = []string{}
	case "createwallet":
		result = map[string]string{"name": "", "warning": ""}
	case "getblockcount":
		result = n.blockCount
	case "getnewaddress":
		result = "bcrt1qfakeaddress"
	case "generatetoaddress":
		hashes := []string{}
		for i := 0; i < int(call.Params[0].(float64)); i++ {
			hashes = append(hashes, fmt.Sprintf("%064x", i))
		}
		result = hashes
	case "rescanblockchain":
		result = map[string]float64{"start_height": 0, "stop_height": n.blockCount}
	case "sendtoaddress":
		result = fmt.Sprintf("%064x", 0xfeed)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"error": map[string]interface{}{"code": -32601, "message": "Method not found"},
		})
		return
	}
	json.NewEncoder(w).Encode(map[string]interface{}{"result": result, "error": nil, "id": 1})
}

func (n *fakeNode) methods() []string {
	n.mu.Lock()
	defer n.mu.Unlock()
	methods := make([]string, 0, len(n.calls))
	for _, c := range n.calls {
		methods = append(methods, c.Method)
	}
	return methods
}

func (n *fakeNode) paramsOf(method string) []interface{} {
	n.mu.Lock()
	defer n.mu.Unlock()
	for _, c := range n.calls {
		if c.Method == method {
			return c.Params
		}
	}
	return nil
}

// fakeConfig is a faucet-enabled Config pointing at a fakeNode.
type fakeConfig struct {
	config.Config
	chain        string
	rpcServerURL string
	registryPath string
}

func (c fakeConfig) IsFaucetEnabled() bool { return true }
func (c fakeConfig) IsMiningEnabled() bool { return true }
func (c fakeConfig) IsLoggerEnabled() bool { return false }
func (c fakeConfig) Chain() string         { return c.chain }
func (c fakeConfig) RPCServerURL() string  { return c.rpcServerURL }
func (c fakeConfig) RegistryPath() string  { return c.registryPath }
func (c fakeConfig) WalletName() string    { return "" }

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}

// Regression test for #60: rescanblockchain must only run on Liquid.
func TestStartupRescanOnlyOnLiquid(t *testing.T) {
	tests := []struct {
		name       string
		chain      string
		blockCount float64
		wantMined  float64
		wantRescan bool
	}{
		{"bitcoin fresh chain", "bitcoin", 0, 101, false},
		{"bitcoin existing chain", "bitcoin", 35000, 0, false},
		{"liquid fresh chain", "liquid", 0, 1, true},
		{"liquid existing chain", "liquid", 35000, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := newFakeNode(t, tt.blockCount)
			r := NewRouter(fakeConfig{
				chain:        tt.chain,
				rpcServerURL: node.URL,
				registryPath: t.TempDir(),
			})
			if r == nil || r.Faucet == nil {
				t.Fatal("Expected *Router with faucet, got <nil>")
			}

			methods := node.methods()

			if !contains(methods, "listwallets") || !contains(methods, "createwallet") {
				t.Fatalf("Expected wallet to be created, got calls: %v", methods)
			}

			if params := node.paramsOf("generatetoaddress"); tt.wantMined > 0 {
				if params == nil || params[0].(float64) != tt.wantMined {
					t.Fatalf("Expected %v blocks mined, got params: %v", tt.wantMined, params)
				}
			} else if params != nil {
				t.Fatalf("Expected no blocks mined on existing chain, got params: %v", params)
			}

			if got := contains(methods, "rescanblockchain"); got != tt.wantRescan {
				t.Fatalf("Expected rescanblockchain called: %v, got: %v (calls: %v)", tt.wantRescan, got, methods)
			}

			resp := faucetRequest(r, btcAddress, 1, "")
			if resp.Code != http.StatusOK {
				t.Fatalf("Expected status %d, got: %d (%s)", http.StatusOK, resp.Code, resp.Body.String())
			}
			if params := node.paramsOf("sendtoaddress"); params == nil || params[0] != btcAddress {
				t.Fatalf("Expected sendtoaddress to %s, got params: %v", btcAddress, params)
			}
		})
	}
}
