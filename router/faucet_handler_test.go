package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

const (
	withLiquid       = true
	btcAddress       = "mpSGWQvbAiRt2UNLST1CdWUufoPVsVwLyK"
	badBtcAddress    = ""
	liquidAddress    = "CTEsqL1x9ooWWG9HBaHUpvS2DGJJ4haYdkTQPKj9U8CCdwT5vcudhbYUT8oQwwoS11aYtdznobfgT8rj"
	badLiquidAddress = ""
	assetQuantity    = 1000
	badAssetQuantity = -1
	name             = "test"
	ticker           = "TST"
)

func TestBitcoinFaucet(t *testing.T) {
	r := NewTestRouter(!withLiquid)

	blockCountResp := blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	prevBlockCount, _ := strconv.Atoi(blockCountResp.Body.String())

	faucetResp := faucetRequest(r, btcAddress)
	if faucetResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, faucetResp.Code)
	}

	// give electrs the time to update block count
	time.Sleep(5 * time.Second)

	blockCountResp = blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	blockCount, _ := strconv.Atoi(blockCountResp.Body.String())
	blockMined := blockCount - prevBlockCount
	if blockMined != 1 {
		t.Fatalf("Expected 1 block mined, got %d\n", blockMined)
	}
}

func TestBitcoinFaucetShouldFail(t *testing.T) {
	r := NewTestRouter(!withLiquid)

	resp := faucetRequest(r, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = faucetRequest(r, badLiquidAddress)
	checkFailed(t, resp, "Invalid address")
}

func TestLiquidFaucet(t *testing.T) {
	r := NewTestRouter(withLiquid)

	blockCountResp := blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	prevBlockCount, _ := strconv.Atoi(blockCountResp.Body.String())

	resp := faucetRequest(r, liquidAddress)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}

	time.Sleep(5 * time.Second)

	blockCountResp = blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	blockCount, _ := strconv.Atoi(blockCountResp.Body.String())
	blockMined := blockCount - prevBlockCount
	if blockMined != 1 {
		t.Fatalf("Expected 1 block mined, got %d\n", blockMined)
	}
	prevBlockCount, _ = strconv.Atoi(blockCountResp.Body.String())

	resp = mintRequest(r, liquidAddress, assetQuantity, name, ticker)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}
	time.Sleep(5 * time.Second)
	blockCountResp = blockCountRequest(r)
	blockCount, _ = strconv.Atoi(blockCountResp.Body.String())
	blockMined = blockCount - prevBlockCount
	if blockMined != 1 {
		t.Fatalf("Expected 1 block mined, got %d\n", blockMined)
	}
}

func TestLiquidFaucetShouldFail(t *testing.T) {
	r := NewTestRouter(withLiquid)

	resp := faucetRequest(r, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = faucetRequest(r, badLiquidAddress)
	checkFailed(t, resp, "Invalid address")

	resp = mintRequest(r, nil, assetQuantity, nil, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = mintRequest(r, liquidAddress, nil, nil, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = mintRequest(r, badLiquidAddress, assetQuantity, nil, nil)
	checkFailed(t, resp, "Invalid address")

	resp = mintRequest(r, liquidAddress, badAssetQuantity, nil, nil)
	checkFailed(t, resp, "Amount out of range")

	resp = mintRequest(r, liquidAddress, assetQuantity, name, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = mintRequest(r, liquidAddress, assetQuantity, nil, name)
	checkFailed(t, resp, "Malformed Request")

	os.RemoveAll(filepath.Join(r.Config.RegistryPath(), "registry"))
}

func faucetRequest(r *Router, address interface{}) *httptest.ResponseRecorder {
	request := map[string]interface{}{
		"address": address,
	}
	payload, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/faucet", bytes.NewBuffer(payload))

	return doRequest(r, req)
}

func blockCountRequest(r *Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/blocks/tip/height", nil)
	return doRequest(r, req)
}

func mintRequest(r *Router, address, quantity, name, ticker interface{}) *httptest.ResponseRecorder {
	request := map[string]interface{}{
		"address":  address,
		"quantity": quantity,
		"name":     name,
		"ticker":   ticker,
	}
	payload, _ := json.Marshal(request)

	req, _ := http.NewRequest("POST", "/mint", bytes.NewBuffer(payload))
	return doRequest(r, req)
}

func doRequest(r *Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}

func checkFailed(t *testing.T, resp *httptest.ResponseRecorder, expectedError string) {
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err := resp.Body.String()
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}
}
