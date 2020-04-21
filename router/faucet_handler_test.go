package router

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
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

func TestBitcoinFaucetBadRequestShouldFail(t *testing.T) {
	r := NewTestRouter(!withLiquid)

	resp := faucetRequest(r, badBtcAddress)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err := resp.Body.String()
	expectedError := "Invalid address"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}

	resp = faucetBadRequest(r)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err = resp.Body.String()
	expectedError = "Malformed Request"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}
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

	resp = mintRequest(r, liquidAddress, assetQuantity)
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

func TestLiquidFaucetBadRequestShouldFail(t *testing.T) {
	r := NewTestRouter(withLiquid)

	resp := faucetRequest(r, badLiquidAddress)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err := resp.Body.String()
	expectedError := "Invalid address"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}

	resp = mintRequest(r, badLiquidAddress, assetQuantity)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err = resp.Body.String()
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}

	resp = mintRequest(r, liquidAddress, badAssetQuantity)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err = resp.Body.String()
	expectedError = "Amount out of range"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}

	resp = faucetBadRequest(r)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err = resp.Body.String()
	expectedError = "Malformed Request"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}

	resp = mintBadRequest(r, "", assetQuantity)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err = resp.Body.String()
	expectedError = "Malformed Request"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}

	resp = mintBadRequest(r, liquidAddress, 0)
	if resp.Code == http.StatusOK {
		t.Fatalf("Should return error, got status: %d\n", resp.Code)
	}

	err = resp.Body.String()
	expectedError = "Malformed Request"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}
}

func faucetRequest(r *Router, address string) *httptest.ResponseRecorder {
	payload := []byte(fmt.Sprintf(`{"address": "%s"}`, address))
	req, _ := http.NewRequest("POST", "/faucet", bytes.NewBuffer(payload))

	return doRequest(r, req)
}

func faucetBadRequest(r *Router) *httptest.ResponseRecorder {
	payload := []byte(fmt.Sprintf("{}"))
	req, _ := http.NewRequest("POST", "/faucet", bytes.NewBuffer(payload))

	return doRequest(r, req)
}

func blockCountRequest(r *Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/blocks/tip/height", nil)
	return doRequest(r, req)
}

func mintRequest(r *Router, address string, quantity float64) *httptest.ResponseRecorder {
	payload := []byte(fmt.Sprintf(`{"address": "%s", "quantity": %f}`, address, quantity))
	req, _ := http.NewRequest("POST", "/mint", bytes.NewBuffer(payload))
	return doRequest(r, req)
}

func mintBadRequest(r *Router, address string, quantity float64) *httptest.ResponseRecorder {
	payload := []byte(fmt.Sprintf(`{"address": "%s", "quantity": %f}`, address, quantity))
	if address == "" {
		payload = []byte(fmt.Sprintf(`{"quantity": %f}`, quantity))
	} else {
		payload = []byte(fmt.Sprintf(`{"address": %s}`, address))
	}
	req, _ := http.NewRequest("POST", "/mint", bytes.NewBuffer(payload))
	return doRequest(r, req)
}

func doRequest(r *Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}
