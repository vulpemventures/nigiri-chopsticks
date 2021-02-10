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
	res.Header().Set("Access-Control-Allow-Origin", "*")

	body := parseRequestBody(req.Body)
	address, ok := body["address"].(string)
	if !ok {
		http.Error(res, "Malformed Request: missing address", http.StatusBadRequest)
		return
	}
	amount, ok := body["amount"].(float64)
	if !ok {
		// the default 100 000 000 satoshis
		amount = 1
	}
	asset, ok := body["asset"].(string)
	if !ok {
		// this means sending bitcoin
		asset = ""
	}

	var status int
	var tx string
	var err error

	if r.Config.Chain() == "liquid" {
		status, tx, err = r.Faucet.SendLiquidTransaction(address, amount, asset)
	} else {
		status, tx, err = r.Faucet.SendBitcoinTransaction(address, amount)
	}
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	if r.Config.IsMiningEnabled() {
		r.Faucet.Mine(1)
	}
	json.NewEncoder(res).Encode(map[string]string{"txId": tx})
	return
}

// HandleMintRequest is a Liquid only endpoint that issues a requested quantity
// of a new asset and sends it to the requested address
func (r *Router) HandleMintRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	body := parseRequestBody(req.Body)
	address, ok := body["address"].(string)
	if !ok {
		http.Error(res, "Malformed Request: missing address", http.StatusBadRequest)
		return
	}
	// NOTICE this is here for backward compatibility. We will deprecate and move to amount
	quantity, qtyOk := body["quantity"].(float64)
	amount, amtOk := body["amount"].(float64)

	if !qtyOk && !amtOk {
		http.Error(res, "Malformed Request: missing amount", http.StatusBadRequest)
		return
	}

	if qtyOk && !amtOk {
		amount = quantity
	}

	status, resp, err := r.Faucet.Mint(address, amount)
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	asset, ok := resp["asset"].(string)
	if !ok {
		http.Error(res, "Internal error", http.StatusInternalServerError)
		return
	}
	issuanceTx, ok := resp["issuance_txin"].(map[string]interface{})
	if !ok {
		http.Error(res, "Internal error", http.StatusInternalServerError)
		return
	}

	name, nameOk := body["name"].(string)
	ticker, tickerOk := body["ticker"].(string)

	if (nameOk && !tickerOk) || (!nameOk && tickerOk) {
		http.Error(res, "Malformed Request: missing name or ticker", http.StatusBadRequest)
		return
	}

	if nameOk && tickerOk {
		contract := map[string]interface{}{
			"name":      name,
			"ticker":    ticker,
			"precision": 8, // we hardcode 8 as precision which is default with issueasset RPC on elements node
		}
		r.Registry.AddEntry(asset, issuanceTx, contract)
	}

	if r.Config.IsMiningEnabled() {
		r.Faucet.Mine(1)
	}

	json.NewEncoder(res).Encode(resp)
	return
}

func parseRequestBody(body io.ReadCloser) map[string]interface{} {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]interface{}
	decoder.Decode(&decodedBody)

	return decodedBody
}
