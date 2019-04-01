package router

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HandleFaucetRequest sends funds to the address passed in the request body
func (r *Router) HandleFaucetRequest(res http.ResponseWriter, req *http.Request) {
	body := parseRequestBody(req.Body)

	status, signedTx, err := r.Faucet.NewTransaction(body["address"])
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	url := fmt.Sprintf("%s/broadcast?tx=%s", r.Config.ElectrsURL(), signedTx)
	status, resp, err := get(url, nil)
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	json.NewEncoder(res).Encode(map[string]string{"txId": resp})
}
