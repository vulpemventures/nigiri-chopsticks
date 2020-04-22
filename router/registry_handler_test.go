package router

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestRegistry(t *testing.T) {
	r := NewTestRouter(withLiquid)
	resp := mintRequest(r, liquidAddress, assetQuantity, name, ticker)
	resp = registryRequest(r, []interface{}{})

	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status %d, got: %d\n", http.StatusOK, resp.Code)
	}

	list := []interface{}{}
	json.Unmarshal(resp.Body.Bytes(), &list)

	if len(list) < 1 {
		t.Fatalf("Expected entry list to not be empty")
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
