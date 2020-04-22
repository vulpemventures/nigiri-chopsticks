package router

import (
	"encoding/json"
	"io"
	"net/http"
)

// HandleFaucetRequest sends funds to the address passed in the request body by
// retrieving the signed transaction from the faucet and broadcasting via the
// electrs broadcast endpoint
func (r *Router) HandleFaucetRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Methods", "POST")

	body := parseRequestBody(req.Body)
	address := body["address"]
	if address == nil {
		http.Error(res, "Malformed Request", http.StatusBadRequest)
		return
	}

	status, tx, err := r.Faucet.NewTransaction(address.(string))
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	if r.Config.IsMiningEnabled() {
		r.Faucet.Mine(1)
		json.NewEncoder(res).Encode(map[string]string{"txId": tx})
	}
	return
}

// HandleMintRequest is a Liquid only endpoint that issues a requested quantity
// of a new asset and sends it to the requested address
func (r *Router) HandleMintRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Methods", "POST")

	body := parseRequestBody(req.Body)
	address := body["address"]
	if address == nil {
		http.Error(res, "Malformed Request", http.StatusBadRequest)
		return
	}
	quantity := body["quantity"]
	if quantity == nil {
		http.Error(res, "Malformed Request", http.StatusBadRequest)
		return
	}
	contract := body["contract"]
	if contract != nil {
		c := contract.(map[string]interface{})
		name := c["name"]
		ticker := c["ticker"]
		if name == nil || ticker == nil {
			http.Error(res, "Malformed Request", http.StatusBadRequest)
			return
		}
	}

	status, resp, err := r.Faucet.Mint(address.(string), quantity.(float64))
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	if contract != nil {
		r.Registry.AddEntry(resp["asset"].(string), resp["issuance_txin"].(map[string]interface{}), contract.(map[string]interface{}))
	}

	if r.Config.IsMiningEnabled() {
		r.Faucet.Mine(1)
		json.NewEncoder(res).Encode(resp)
	}
	return
}

func parseRequestBody(body io.ReadCloser) map[string]interface{} {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]interface{}
	decoder.Decode(&decodedBody)

	return decodedBody
}
