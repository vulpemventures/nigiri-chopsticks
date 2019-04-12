package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const (
	address       = "mpSGWQvbAiRt2UNLST1CdWUufoPVsVwLyK"
	liquidAddress = "CTEsqL1x9ooWWG9HBaHUpvS2DGJJ4haYdkTQPKj9U8CCdwT5vcudhbYUT8oQwwoS11aYtdznobfgT8rj"

	badAddress       = "mN9Zw9TwwKoYXrBzVFuioXab3xeu4WWr1wk"
	badLiquidAddress = "CTEqd9wVbwWwNnXuDbN9q6NMWCpHX6FeqaYgBsafvpQhQV2v5aHWGkYrns43ogncYcjksru2cLVa7MJv"
)

// For bitcoin, this implicitly tests also the broadcast endpoint since
// for this chain the faucet gets a signed transaction from the faucet
// and broadcasts it via the cited endpoint
func TestFaucetHandler(t *testing.T) {
	rr := makeFaucetRequest(false, address)
	resp := rr.Body.String()
	status := rr.Code
	expectedStatus := http.StatusOK

	if status != expectedStatus {
		t.Fatal(resp)
	}

	t.Log("tx hash:", resp)
}

func TestFaucetHandlerShouldFail(t *testing.T) {
	rr := makeFaucetRequest(false, badAddress)
	resp := strings.TrimSpace(rr.Body.String())
	status := rr.Code
	expectedStatus := http.StatusInternalServerError
	expectedResp := fmt.Sprintf("Method createrawtransaction failed with error: Invalid Bitcoin address: %s", badAddress)

	if status != expectedStatus {
		t.Fatalf("Expected status: %d, got: %d\n check out the response: %s", expectedStatus, status, resp)
	}

	if resp != expectedResp {
		t.Fatalf("Expected response: %s, got: %s", expectedResp, resp)
	}
}

func TestLiquidFaucetHandler(t *testing.T) {
	rr := makeFaucetRequest(true, liquidAddress)
	resp := rr.Body.String()
	status := rr.Code
	expectedStatus := http.StatusOK

	if status != expectedStatus {
		t.Fatal(resp)
	}

	t.Log("tx hash:", resp)
}

func TestLiquidFaucetHandlerShouldFail(t *testing.T) {
	rr := makeFaucetRequest(true, badLiquidAddress)
	resp := strings.TrimSpace(rr.Body.String())
	status := rr.Code
	expectedStatus := http.StatusInternalServerError
	expectedResp := "Method sendtoaddress failed with error: Invalid Bitcoin address"

	if status != expectedStatus {
		t.Fatalf("Expected status: %d, got: %d\n check out the response: %s", expectedStatus, status, resp)
	}

	if resp != expectedResp {
		t.Fatalf("Expected response: %s, got: %s", expectedResp, resp)
	}
}

func makeFaucetRequest(isLiquid bool, addr string) *httptest.ResponseRecorder {
	router := newTestingRouter(isLiquid)

	payloadBuffer := &bytes.Buffer{}
	json.NewEncoder(payloadBuffer).Encode(map[string]string{"address": addr})

	req, _ := http.NewRequest("POST", "/faucet", payloadBuffer)
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(router.HandleFaucetRequest)

	handler.ServeHTTP(rr, req)

	return rr
}
