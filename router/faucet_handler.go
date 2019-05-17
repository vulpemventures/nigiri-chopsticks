package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// HandleFaucetRequest sends funds to the address passed in the request body by
// retrieving the signed transaction from the faucet and broadcasting via the
// electrs broadcast endpoint
func (r *Router) HandleFaucetRequest(res http.ResponseWriter, req *http.Request) {
	body := parseRequestBody(req.Body)

	status, tx, err := r.Faucet.NewTransaction(body["address"])
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	/*
		since the liquid implementation of Faucet uses the sendtoaddress RPC method,
		it already returns the hash of the transaction, then we just need to get it
		confirmed mining one block and returning the tx hash in the response body
	*/
	if r.Config.Chain() == "liquid" {
		r.Faucet.Mine(10)
		json.NewEncoder(res).Encode(map[string]string{"txId": tx})
		return
	}

	broadcastRequest, _ := http.NewRequest("GET", fmt.Sprintf("%s/broadcast?tx=%s", r.Config.ElectrsURL(), tx), nil)

	r.HandleBroadcastRequest(res, broadcastRequest)
}

func parseRequestBody(body io.ReadCloser) map[string]string {
	decoder := json.NewDecoder(body)
	var decodedBody map[string]string
	decoder.Decode(&decodedBody)

	return decodedBody
}
