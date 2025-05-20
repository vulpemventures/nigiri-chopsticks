package router

import (
	"encoding/json"
	"net/http"

	"github.com/vulpemventures/nigiri-chopsticks/helpers"
)

// HandleSubmitPackageRequest handles the submitpackage request
// it mocks the incoming /txs/package endpoint
// TODO remove once esplora supports it
func (r *Router) HandleSubmitPackageRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	decoder := json.NewDecoder(req.Body)
	var txs []string
	err := decoder.Decode(&txs)
	if err != nil {
		http.Error(res, "Malformed Request: missing txs", http.StatusBadRequest)
		return
	}

	status, resp, err := helpers.HandleRPCRequest(r.RPCClient, "submitpackage", []interface{}{txs})
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	respMap, ok := resp.(map[string]interface{})
	if !ok {
		http.Error(res, "Malformed Response: expected JSON object", http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(res).Encode(respMap); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
}
