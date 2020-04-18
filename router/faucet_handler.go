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

	status, tx, err := r.Faucet.NewTransaction(body["address"].(string))
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
	address := body["address"].(string)
	quantity := body["quantity"].(float64)

	status, resp, err := r.Faucet.Mint(address, quantity)
	if err != nil {
		http.Error(res, err.Error(), status)
		return
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
