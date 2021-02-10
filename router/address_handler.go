package router

import (
	"encoding/json"
	"net/http"
)

// HandleAddressRequest calls the `getnewaddress` and returns the native segwit one
func (r *Router) HandleAddressRequest(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Access-Control-Allow-Origin", "*")

	status, resp, err := r.RPCClient.Call("getnewaddress", []interface{}{"", "bech32"})
	if err != nil {
		http.Error(res, err.Error(), status)
		return
	}

	var out interface{}
	err = json.Unmarshal(resp.Result, &out)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	addr := out.(string)

	json.NewEncoder(res).Encode(map[string]string{"address": addr})
	return
}
