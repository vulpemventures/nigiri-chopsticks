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
	expectedError := "Invalid Bitcoin address"
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
	if blockMined != 10 {
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
	expectedError := "Invalid Bitcoin address"
	if !strings.Contains(err, expectedError) {
		t.Fatalf("Expected error: %s, got: %s\n", expectedError, err)
	}
}

func faucetRequest(r *Router, address string) *httptest.ResponseRecorder {
	payload := []byte(fmt.Sprintf(`{"address": "%s"}`, address))
	req, _ := http.NewRequest("POST", "/faucet", bytes.NewBuffer(payload))

	return doRequest(r, req)
}

func blockCountRequest(r *Router) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", "/blocks/tip/height", nil)
	return doRequest(r, req)
}

func doRequest(r *Router, req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)
	return rr
}
