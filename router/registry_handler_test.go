package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRegistry(t *testing.T) {
	r := NewTestRouter(withLiquid)
	resp := mintRequest(r, liquidAddress, assetQuantity, name, ticker)

	parsedResp := map[string]interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &parsedResp)
	asset := parsedResp["asset"].(string)

	resp = registryRequest(r, []interface{}{})
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}

	list := []map[string]interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &list)

	if len(list) < 1 {
		t.Fatalf("Expected entry list to not be empty")
	}

	resAsset := list[0]["asset"].(string)
	resName := list[0]["name"].(string)
	resTicker := list[0]["ticker"].(string)

	if strings.Compare(asset, resAsset) != 0 {
		t.Fatalf("Expected asset hash: %s, got: %s", asset, resAsset)
	}
	if strings.Compare(name, resName) != 0 {
		t.Fatalf("Expected name hash: %s, got: %s", name, resName)
	}
	if strings.Compare(ticker, resTicker) != 0 {
		t.Fatalf("Expected ticker hash: %s, got: %s", ticker, resTicker)
	}
}

func TestRegistryShouldFail(t *testing.T) {
	r := NewTestRouter(withLiquid)
	resp := registryRequest(r, nil)
	checkFailed(t, resp, "Malformed Request")

	os.RemoveAll(filepath.Join(r.Config.RegistryPath(), "registry"))
}

func registryRequest(r *Router, assets []interface{}) *httptest.ResponseRecorder {
	request := map[string]interface{}{
		"assets": assets,
	}

	payload, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/registry", bytes.NewBuffer(payload))
	return doRequest(r, req)
}
