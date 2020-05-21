package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestAssetEndpointWithExtraInfo(t *testing.T) {
	r := NewTestRouter(withLiquid)
	resp := mintRequest(r, liquidAddress, assetQuantity, name, ticker)

	parsedResp := map[string]interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &parsedResp)
	asset := parsedResp["asset"].(string)

	time.Sleep(5 * time.Second)

	resp = assetRequest(r, asset)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}

	parsedResp = map[string]interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &parsedResp)
	resName := parsedResp["name"].(string)
	resTicker := parsedResp["ticker"].(string)

	if strings.Compare(name, resName) != 0 {
		t.Fatalf("Expected asset name: %s, got: %s\n", name, resName)
	}
	if strings.Compare(ticker, resTicker) != 0 {
		t.Fatalf("Expected asset ticker: %s, got: %s\n", ticker, resTicker)
	}
}

func TestAssetEndpointWithoutExtraInfo(t *testing.T) {
	r := NewTestRouter(withLiquid)
	resp := mintRequest(r, liquidAddress, assetQuantity, nil, nil)

	parsedResp := map[string]interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &parsedResp)
	asset := parsedResp["asset"].(string)

	time.Sleep(5 * time.Second)

	resp = assetRequest(r, asset)
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}

	parsedResp = map[string]interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &parsedResp)

	if parsedResp["name"] != nil {
		t.Fatalf("Asset name should not be defined")
	}

	if parsedResp["ticker"] != nil {
		t.Fatalf("Asset ticker should not be defined")
	}
}

func TestAssetEndpointForNonExistingAssetID(t *testing.T) {
	r := NewTestRouter(withLiquid)
	assetID := "dummyID"
	resp := assetRequest(r, assetID)

	if resp.Code != http.StatusBadRequest {
		t.Errorf("for non existing assetID: %v, expected status code is 400, actual: %v\n", assetID, resp.Code)
	}
}

func assetRequest(r *Router, asset interface{}) *httptest.ResponseRecorder {
	path := fmt.Sprintf("/asset/%s", asset)
	req, _ := http.NewRequest("GET", path, nil)
	return doRequest(r, req)
}
