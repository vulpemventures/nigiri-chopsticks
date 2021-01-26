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

	faucetResp := faucetRequest(r, btcAddress, 0, "")
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

	resp := faucetRequest(r, "", 0, "")
	checkFailed(t, resp, "Invalid address")

	resp = faucetRequest(r, badLiquidAddress, 0, "")
	checkFailed(t, resp, "Invalid address")
}

func TestLiquidFaucet(t *testing.T) {
	r := NewTestRouter(withLiquid)

	blockCountResp := blockCountRequest(r)
	if blockCountResp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, blockCountResp.Code)
	}
	prevBlockCount, _ := strconv.Atoi(blockCountResp.Body.String())

	resp := faucetRequest(r, liquidAddress, 0, "")
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

	addrOfNode, err := getNewAddressFromNode(r)
	if err != nil {
		t.Fatalf("getNewAddressFromNode error")
	}

	resp = mintRequest(r, addrOfNode, assetQuantity, name, ticker)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}
	time.Sleep(5 * time.Second)

	var decodedBody map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &decodedBody)
	if err != nil {
		t.Fatalf("decoding response error")
	}

	resp2 := faucetRequest(r, liquidAddress, 0, decodedBody["asset"].(string))
	if resp2.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp2.Code)
	}
	time.Sleep(5 * time.Second)
	blockCountResp = blockCountRequest(r)
	blockCount, _ = strconv.Atoi(blockCountResp.Body.String())
	blockMined = blockCount - prevBlockCount
	if blockMined < 1 {
		t.Fatalf("Expected at least 1 block mined, got %d\n", blockMined)
	}

}

func TestLiquidFaucetShouldFail(t *testing.T) {
	r := NewTestRouter(withLiquid)

	resp := faucetRequest(r, "", 0, "")
	checkFailed(t, resp, "Invalid address")

	resp = faucetRequest(r, badLiquidAddress, 0, "")
	checkFailed(t, resp, "Invalid address")

	resp = mintRequest(r, "", assetQuantity, nil, nil)
	checkFailed(t, resp, "Invalid address")

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

func faucetRequest(r *Router, address string, amount float64, asset string) *httptest.ResponseRecorder {

	request := map[string]interface{}{
		"address": address,
	}

	if amount > 0 {
		request["amount"] = amount
	}

	if len(asset) > 0 {
		request["asset"] = asset
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

func getNewAddressFromNode(r *Router) (string, error) {
	_, resp, err := handleRPCRequest(r.RPCClient, "getnewaddress", nil)
	if err != nil {
		return "", err
	}

	return resp.(string), nil
}
