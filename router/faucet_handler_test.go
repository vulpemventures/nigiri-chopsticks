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
)

var (
	contract = map[string]string{
		"name":   "test",
		"ticker": "TST",
	}
	badContract1 = map[string]string{
		"name": "test",
	}
	badContract2 = map[string]string{
		"ticker": "TST",
	}
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

	resp = mintRequest(r, liquidAddress, assetQuantity, contract)
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

	resp = mintRequest(r, nil, assetQuantity, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = mintRequest(r, liquidAddress, nil, nil)
	checkFailed(t, resp, "Malformed Request")

	resp = mintRequest(r, badLiquidAddress, assetQuantity, nil)
	checkFailed(t, resp, "Invalid address")

	resp = mintRequest(r, liquidAddress, badAssetQuantity, nil)
	checkFailed(t, resp, "Amount out of range")

	resp = mintRequest(r, liquidAddress, assetQuantity, badContract1)
	checkFailed(t, resp, "Malformed Request")

	resp = mintRequest(r, liquidAddress, assetQuantity, badContract2)
	checkFailed(t, resp, "Malformed Request")

	os.Remove(filepath.Join(r.Config.RegistryPath(), "registry.json"))
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

func mintRequest(r *Router, address, quantity, contract interface{}) *httptest.ResponseRecorder {
	request := map[string]interface{}{
		"address":  address,
		"quantity": quantity,
		"contract": contract,
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
